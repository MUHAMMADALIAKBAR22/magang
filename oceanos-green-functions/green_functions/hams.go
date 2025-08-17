// Package green_functions - HAMS and LiangWuNoblesse implementations
// Copyright (C) 2025 Capytaine Contributors
// See LICENSE file at <https://github.com/capytaine/capytaine>

package green_functions

import (
	"errors"
	"gonum.org/v1/gonum/mat"
	"math"
)

// LiangWuNoblesseGF implements the infinite depth Green function of Liang, Wu, Noblesse (2018)
// Uses the same implementation as Delhommeau for the Rankine and reflected Rankine terms
type LiangWuNoblesseGF struct {
	*BaseGreenFunction
	exportableSettings map[string]interface{}
}

// NewLiangWuNoblesseGF creates a new LiangWuNoblesse Green function
func NewLiangWuNoblesseGF() *LiangWuNoblesseGF {
	lwn := &LiangWuNoblesseGF{
		BaseGreenFunction: NewBaseGreenFunction(),
		exportableSettings: map[string]interface{}{
			"green_function": "LiangWuNoblesseGF",
		},
	}
	lwn.SetFloatingPointPrecision(Float64)
	return lwn
}

// String returns a string representation of the LiangWuNoblesse Green function
func (lwn *LiangWuNoblesseGF) String() string {
	return "LiangWuNoblesseGF()"
}

// Evaluate computes the Green function using the LiangWuNoblesse method
func (lwn *LiangWuNoblesseGF) Evaluate(mesh1, mesh2 interface{}, freeSurface float64, waterDepth float64,
	wavenumber complex128, adjointDoubleLayer bool, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error) {

	// Check constraints for LiangWuNoblesse method
	if math.IsInf(freeSurface, 1) || !math.IsInf(waterDepth, 1) {
		return nil, nil, errors.New("LiangWuNoblesseGF is only implemented for infinite depth with a free surface")
	}

	// Get collocation points and normals
	colocationPoints, earlyDotProductNormals, err := lwn.getColocationPointsAndNormals(mesh1, mesh2, adjointDoubleLayer)
	if err != nil {
		return nil, nil, err
	}

	// Determine matrix dimensions
	var rows, cols int
	if meshLike1, ok := mesh1.(MeshLike); ok {
		rows = meshLike1.GetNbFaces()
	} else if pointArray, ok := mesh1.(*mat.Dense); ok {
		rows, _ = pointArray.Dims()
	}

	if meshLike2, ok := mesh2.(MeshLike); ok {
		cols = meshLike2.GetNbFaces()
	} else {
		return nil, nil, &GreenFunctionEvaluationError{"mesh2 must implement MeshLike interface"}
	}

	// Initialize matrices
	S, K, err := lwn.initMatrices(rows, cols, earlyDotProduct)
	if err != nil {
		return nil, nil, err
	}

	// Determine singularity handling based on wavenumber
	var gfSingularitiesIndex int
	if math.IsInf(real(wavenumber), 1) {
		gfSingularitiesIndex = 0 // high_freq
	} else {
		gfSingularitiesIndex = 1 // low_freq
	}

	// TODO: Implement actual LiangWuNoblesse Green function computation
	// This would involve:
	// 1. Computing the Rankine and reflected Rankine terms
	// 2. Applying the LiangWuNoblesse method for wave terms
	// 3. Handling singularities appropriately

	_ = colocationPoints
	_ = earlyDotProductNormals
	_ = gfSingularitiesIndex

	return S, K, nil
}

// HAMS represents the HAMS (Hydrodynamic Analysis of Marine Structures) Green function
type HAMS struct {
	*BaseGreenFunction
	exportableSettings map[string]interface{}
}

// NewHAMS creates a new HAMS Green function
func NewHAMS() *HAMS {
	hams := &HAMS{
		BaseGreenFunction: NewBaseGreenFunction(),
		exportableSettings: map[string]interface{}{
			"green_function": "HAMS",
		},
	}
	hams.SetFloatingPointPrecision(Float64)
	return hams
}

// String returns a string representation of the HAMS Green function
func (h *HAMS) String() string {
	return "HAMS()"
}

// Evaluate computes the Green function using the HAMS method
func (h *HAMS) Evaluate(mesh1, mesh2 interface{}, freeSurface float64, waterDepth float64,
	wavenumber complex128, adjointDoubleLayer bool, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error) {

	// Get collocation points and normals
	colocationPoints, earlyDotProductNormals, err := h.getColocationPointsAndNormals(mesh1, mesh2, adjointDoubleLayer)
	if err != nil {
		return nil, nil, err
	}

	// Determine matrix dimensions
	var rows, cols int
	if meshLike1, ok := mesh1.(MeshLike); ok {
		rows = meshLike1.GetNbFaces()
	} else if pointArray, ok := mesh1.(*mat.Dense); ok {
		rows, _ = pointArray.Dims()
	}

	if meshLike2, ok := mesh2.(MeshLike); ok {
		cols = meshLike2.GetNbFaces()
	} else {
		return nil, nil, &GreenFunctionEvaluationError{"mesh2 must implement MeshLike interface"}
	}

	// Initialize matrices
	S, K, err := h.initMatrices(rows, cols, earlyDotProduct)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Implement actual HAMS Green function computation
	// This would involve the HAMS-specific method for computing Green functions

	_ = colocationPoints
	_ = earlyDotProductNormals
	_ = freeSurface
	_ = waterDepth
	_ = wavenumber

	return S, K, nil
}
