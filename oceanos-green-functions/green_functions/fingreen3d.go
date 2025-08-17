// Package green_functions - FinGreen3D implementation for finite depth Green functions
// Copyright (C) 2025 Capytaine Contributors
// See LICENSE file at <https://github.com/capytaine/capytaine>

package green_functions

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"math"
	"math/cmplx"
)

// FinGreen3D implements finite depth Green function computation
// Based on the Fortran implementation by Yingyi Liu (2013)
type FinGreen3D struct {
	*BaseGreenFunction
	waterDepth         float64
	waveNumber         complex128
	dispersionRoots    []complex128
	exportableSettings map[string]interface{}
}

// NewFinGreen3D creates a new FinGreen3D Green function
func NewFinGreen3D(waterDepth float64) *FinGreen3D {
	fg := &FinGreen3D{
		BaseGreenFunction: NewBaseGreenFunction(),
		waterDepth:        waterDepth,
		exportableSettings: map[string]interface{}{
			"green_function": "FinGreen3D",
			"water_depth":    waterDepth,
		},
	}
	fg.SetFloatingPointPrecision(Float64)
	return fg
}

// String returns a string representation of the FinGreen3D Green function
func (fg *FinGreen3D) String() string {
	return "FinGreen3D(water_depth=" + fmt.Sprintf("%.2f", fg.waterDepth) + ")"
}

// SetWaveNumber sets the wave number and computes dispersion relation roots
func (fg *FinGreen3D) SetWaveNumber(wavenumber complex128) error {
	fg.waveNumber = wavenumber

	// Compute dispersion relation roots
	// For finite depth: omega^2 = g*k*tanh(k*h)
	// where h is water depth
	roots, err := fg.computeDispersionRoots(wavenumber)
	if err != nil {
		return err
	}
	fg.dispersionRoots = roots

	return nil
}

// computeDispersionRoots computes the roots of the dispersion relation
func (fg *FinGreen3D) computeDispersionRoots(wavenumber complex128) ([]complex128, error) {
	// TODO: Implement proper dispersion relation solver
	// This involves solving: V = k*tanh(k*h) where V = omega^2/g

	// For now, return a simple approximation
	k0 := wavenumber
	roots := []complex128{k0}

	// Add imaginary roots for finite depth case
	if !math.IsInf(fg.waterDepth, 1) {
		// Add additional roots for finite depth
		for n := 1; n <= 10; n++ {
			kn := complex(0, float64(n)*math.Pi/fg.waterDepth)
			roots = append(roots, kn)
		}
	}

	return roots, nil
}

// Evaluate computes the Green function using FinGreen3D method
func (fg *FinGreen3D) Evaluate(mesh1, mesh2 interface{}, freeSurface float64, waterDepth float64,
	wavenumber complex128, adjointDoubleLayer bool, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error) {

	// Update water depth and wave number if changed
	if waterDepth != fg.waterDepth {
		fg.waterDepth = waterDepth
	}

	err := fg.SetWaveNumber(wavenumber)
	if err != nil {
		return nil, nil, err
	}

	// Get collocation points and normals
	colocationPoints, earlyDotProductNormals, err := fg.getColocationPointsAndNormals(mesh1, mesh2, adjointDoubleLayer)
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
	S, K, err := fg.initMatrices(rows, cols, earlyDotProduct)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Implement actual FinGreen3D computation
	// This would involve:
	// 1. Computing the Rankine part: 1/r + 1/r'
	// 2. Computing the wave part Gw using series expansion
	// 3. Handling singularities properly

	_ = colocationPoints
	_ = earlyDotProductNormals
	_ = freeSurface

	return S, K, nil
}

// computeGreenFunction3D computes the 3D finite depth Green function
func (fg *FinGreen3D) computeGreenFunction3D(rr, zf, zp float64) (complex128, error) {
	// Input parameters:
	// rr: horizontal distance between field and source point
	// zf: z coordinate of field point
	// zp: z coordinate of source point

	if math.IsInf(fg.waterDepth, 1) {
		// Infinite depth case
		return fg.computeInfiniteDepthGF(rr, zf, zp)
	}

	// Finite depth case
	return fg.computeFiniteDepthGF(rr, zf, zp)
}

// computeInfiniteDepthGF computes infinite depth Green function
func (fg *FinGreen3D) computeInfiniteDepthGF(rr, zf, zp float64) (complex128, error) {
	// Rankine parts
	r := math.Sqrt(rr*rr + (zf-zp)*(zf-zp))
	rPrime := math.Sqrt(rr*rr + (zf+zp)*(zf+zp))

	rankine := 1.0/r + 1.0/rPrime

	// Wave part for infinite depth
	k := real(fg.waveNumber)
	if k == 0 {
		return complex(rankine, 0), nil
	}

	// Simplified wave part computation
	waveHeight := k * (zf + zp)
	wavePart := complex(0, 0)

	if waveHeight > -10 { // Avoid numerical issues
		wavePart = complex(2*k*math.Exp(waveHeight), 0)
	}

	return complex(rankine, 0) + wavePart, nil
}

// computeFiniteDepthGF computes finite depth Green function
func (fg *FinGreen3D) computeFiniteDepthGF(rr, zf, zp float64) (complex128, error) {
	// Rankine parts
	r := math.Sqrt(rr*rr + (zf-zp)*(zf-zp))
	rPrime := math.Sqrt(rr*rr + (zf+zp)*(zf+zp))

	rankine := 1.0/r + 1.0/rPrime

	// Wave part using dispersion roots
	var wavePart complex128 = 0

	for i, root := range fg.dispersionRoots {
		if i == 0 {
			// First root (propagating mode)
			k := root
			if real(k) > 0 {
				term := fg.computeWaveTerm(rr, zf, zp, k, true)
				wavePart += term
			}
		} else {
			// Higher order roots (evanescent modes)
			kn := root
			term := fg.computeWaveTerm(rr, zf, zp, kn, false)
			wavePart += term
		}
	}

	return complex(rankine, 0) + wavePart, nil
}

// computeWaveTerm computes individual wave terms
func (fg *FinGreen3D) computeWaveTerm(rr, zf, zp float64, k complex128, isPropagating bool) complex128 {
	h := fg.waterDepth

	// Vertical functions
	var phi1, phi2 complex128

	if isPropagating {
		// Propagating mode
		coshKZ1 := cmplx.Cosh(k * complex(zf+h, 0))
		coshKZ2 := cmplx.Cosh(k * complex(zp+h, 0))
		coshKH := cmplx.Cosh(k * complex(h, 0))

		phi1 = coshKZ1 / coshKH
		phi2 = coshKZ2 / coshKH
	} else {
		// Evanescent mode
		cosKZ1 := cmplx.Cos(k * complex(zf+h, 0))
		cosKZ2 := cmplx.Cos(k * complex(zp+h, 0))
		cosKH := cmplx.Cos(k * complex(h, 0))

		phi1 = cosKZ1 / cosKH
		phi2 = cosKZ2 / cosKH
	}

	// Horizontal function (simplified Bessel function approximation)
	kr := k * complex(rr, 0)
	var horizontalFunc complex128

	if cmplx.Abs(kr) < 0.1 {
		// Small argument approximation
		horizontalFunc = complex(1, 0)
	} else {
		// Large argument approximation
		horizontalFunc = cmplx.Exp(1i*kr) / cmplx.Sqrt(kr)
	}

	return phi1 * phi2 * horizontalFunc
}
