// Package green_functions provides Green function computation methods for boundary element method
// Copyright (C) 2025 Capytaine Contributors
// See LICENSE file at <https://github.com/capytaine/capytaine>

package green_functions

import (
	"errors"
	"gonum.org/v1/gonum/mat"
)

// GreenFunctionEvaluationError represents errors during Green function evaluation
type GreenFunctionEvaluationError struct {
	Message string
}

func (e *GreenFunctionEvaluationError) Error() string {
	return e.Message
}

// FloatingPointPrecision represents the precision for floating point calculations
type FloatingPointPrecision string

const (
	Float32 FloatingPointPrecision = "float32"
	Float64 FloatingPointPrecision = "float64"
)

// MeshLike represents a mesh-like structure with faces and normals
type MeshLike interface {
	GetFacesCenters() *mat.Dense
	GetFacesNormals() *mat.Dense
	GetNbFaces() int
}

// AbstractGreenFunction defines the interface for Green function implementations
type AbstractGreenFunction interface {
	// Evaluate computes the Green function between two meshes
	Evaluate(mesh1, mesh2 interface{}, freeSurface float64, waterDepth float64,
		wavenumber complex128, adjointDoubleLayer bool, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error)

	// GetFloatingPointPrecision returns the precision used for calculations
	GetFloatingPointPrecision() FloatingPointPrecision

	// SetFloatingPointPrecision sets the precision for calculations
	SetFloatingPointPrecision(precision FloatingPointPrecision)
}

// BaseGreenFunction provides common functionality for Green function implementations
type BaseGreenFunction struct {
	FloatingPointPrecision FloatingPointPrecision
}

// NewBaseGreenFunction creates a new base Green function with default precision
func NewBaseGreenFunction() *BaseGreenFunction {
	return &BaseGreenFunction{
		FloatingPointPrecision: Float64,
	}
}

// GetFloatingPointPrecision returns the current floating point precision
func (bgf *BaseGreenFunction) GetFloatingPointPrecision() FloatingPointPrecision {
	return bgf.FloatingPointPrecision
}

// SetFloatingPointPrecision sets the floating point precision
func (bgf *BaseGreenFunction) SetFloatingPointPrecision(precision FloatingPointPrecision) {
	bgf.FloatingPointPrecision = precision
}

// getColocationPointsAndNormals extracts collocation points and normals from mesh inputs
func (bgf *BaseGreenFunction) getColocationPointsAndNormals(mesh1, mesh2 interface{}, adjointDoubleLayer bool) (*mat.Dense, *mat.Dense, error) {
	var colocationPoints *mat.Dense
	var earlyDotProductNormals *mat.Dense

	switch m1 := mesh1.(type) {
	case MeshLike:
		colocationPoints = m1.GetFacesCenters()
		if !adjointDoubleLayer {
			// Computing the D matrix
			if m2, ok := mesh2.(MeshLike); ok {
				earlyDotProductNormals = m2.GetFacesNormals()
			} else {
				return nil, nil, errors.New("mesh2 must implement MeshLike interface when mesh1 is MeshLike")
			}
		} else {
			// Computing the K matrix
			earlyDotProductNormals = m1.GetFacesNormals()
		}
	case *mat.Dense:
		// This is used when computing potential or velocity at given points in postprocessing
		rows, cols := m1.Dims()
		if cols != 3 {
			return nil, nil, errors.New("point array must have 3 columns (x, y, z coordinates)")
		}
		colocationPoints = m1

		if !adjointDoubleLayer {
			// Computing the D matrix
			if m2, ok := mesh2.(MeshLike); ok {
				earlyDotProductNormals = m2.GetFacesNormals()
			} else {
				return nil, nil, errors.New("mesh2 must implement MeshLike interface when mesh1 is point array")
			}
		} else {
			// Computing the K matrix - dummy normals
			earlyDotProductNormals = mat.NewDense(rows, 3, nil)
		}
	default:
		return nil, nil, errors.New("unrecognized first input type for Green function evaluation")
	}

	return colocationPoints, earlyDotProductNormals, nil
}

// initMatrices initializes the S and K matrices with appropriate dimensions and data types
func (bgf *BaseGreenFunction) initMatrices(rows, cols int, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error) {
	// Initialize S matrix (single layer potential)
	S := mat.NewCDense(rows, cols, nil)

	// Initialize K matrix (double layer potential)
	var kCols int
	if earlyDotProduct {
		kCols = 1
	} else {
		kCols = 3
	}
	K := mat.NewCDense(rows, cols*kCols, nil)

	return S, K, nil
}
