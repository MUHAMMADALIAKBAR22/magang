// Package green_functions - Delhommeau method implementation
// Copyright (C) 2025 Capytaine Contributors
// See LICENSE file at <https://github.com/capytaine/capytaine>

package green_functions

import (
	"fmt"
	"gonum.org/v1/gonum/mat"
	"hash/fnv"
	"path/filepath"
	"sort"
)

// TabulationGridShape represents different grid shapes for tabulation
type TabulationGridShape string

const (
	Legacy       TabulationGridShape = "legacy"
	ScaledNemoh3 TabulationGridShape = "scaled_nemoh3"
)

// FiniteDepthMethod represents different methods for finite depth calculations
type FiniteDepthMethod string

const (
	LegacyMethod FiniteDepthMethod = "legacy"
	NewerMethod  FiniteDepthMethod = "newer"
)

// GFSingularities represents different singularity handling methods
type GFSingularities string

const (
	HighFreq               GFSingularities = "high_freq"
	LowFreq                GFSingularities = "low_freq"
	LowFreqWithRankinePart GFSingularities = "low_freq_with_rankine_part"
)

// PronyDecompositionMethod represents different Prony decomposition methods
type PronyDecompositionMethod string

const (
	PythonMethod  PronyDecompositionMethod = "python"
	FortranMethod PronyDecompositionMethod = "fortran"
)

// DelhommeauParameters holds configuration parameters for Delhommeau method
type DelhommeauParameters struct {
	TabulationNr                        int
	TabulationRmax                      float64
	TabulationNz                        int
	TabulationZmin                      float64
	TabulationNbIntegrationPoints       int
	TabulationGridShape                 TabulationGridShape
	TabulationCacheDir                  string
	FiniteDepthMethod                   FiniteDepthMethod
	FiniteDepthPronyDecompositionMethod PronyDecompositionMethod
	FloatingPointPrecision              FloatingPointPrecision
	GfSingularities                     GFSingularities
}

// DefaultDelhommeauParameters returns default parameters for Delhommeau method
func DefaultDelhommeauParameters() DelhommeauParameters {
	return DelhommeauParameters{
		TabulationNr:                        676,
		TabulationRmax:                      100.0,
		TabulationNz:                        372,
		TabulationZmin:                      -251.0,
		TabulationNbIntegrationPoints:       1001,
		TabulationGridShape:                 ScaledNemoh3,
		FiniteDepthMethod:                   NewerMethod,
		FiniteDepthPronyDecompositionMethod: PythonMethod,
		FloatingPointPrecision:              Float64,
		GfSingularities:                     LowFreq,
	}
}

// Delhommeau implements the Green function as in Aquadyn and Nemoh
type Delhommeau struct {
	*BaseGreenFunction
	parameters               DelhommeauParameters
	tabulationGridShapeIndex int
	finiteDepthMethodIndex   int
	gfSingularitiesIndex     int
	dispersionRelationRoots  []complex128
	exportableSettings       map[string]interface{}
	hash                     uint64
}

// NewDelhommeau creates a new Delhommeau Green function with specified parameters
func NewDelhommeau(params DelhommeauParameters) *Delhommeau {
	d := &Delhommeau{
		BaseGreenFunction:       NewBaseGreenFunction(),
		parameters:              params,
		dispersionRelationRoots: make([]complex128, 1), // dummy array
	}

	d.SetFloatingPointPrecision(params.FloatingPointPrecision)

	// Set grid shape index
	switch params.TabulationGridShape {
	case Legacy:
		d.tabulationGridShapeIndex = 0
	case ScaledNemoh3:
		d.tabulationGridShapeIndex = 1
	}

	// Set finite depth method index
	switch params.FiniteDepthMethod {
	case LegacyMethod:
		d.finiteDepthMethodIndex = 0
	case NewerMethod:
		d.finiteDepthMethodIndex = 1
	}

	// Set GF singularities index
	switch params.GfSingularities {
	case HighFreq:
		d.gfSingularitiesIndex = 0
	case LowFreq:
		d.gfSingularitiesIndex = 1
	case LowFreqWithRankinePart:
		d.gfSingularitiesIndex = 2
	}

	// Create exportable settings
	d.exportableSettings = map[string]interface{}{
		"green_function":                          "Delhommeau",
		"tabulation_nr":                           params.TabulationNr,
		"tabulation_rmax":                         params.TabulationRmax,
		"tabulation_nz":                           params.TabulationNz,
		"tabulation_zmin":                         params.TabulationZmin,
		"tabulation_nb_integration_points":        params.TabulationNbIntegrationPoints,
		"tabulation_grid_shape":                   params.TabulationGridShape,
		"finite_depth_method":                     params.FiniteDepthMethod,
		"finite_depth_prony_decomposition_method": params.FiniteDepthPronyDecompositionMethod,
		"floating_point_precision":                params.FloatingPointPrecision,
		"gf_singularities":                        params.GfSingularities,
	}

	d.hash = d.computeHash()

	// Initialize tabulation
	if params.TabulationCacheDir == "" {
		d.createTabulation()
	} else {
		d.createOrLoadTabulation()
	}

	return d
}

// NewDefaultDelhommeau creates a new Delhommeau Green function with default parameters
func NewDefaultDelhommeau() *Delhommeau {
	return NewDelhommeau(DefaultDelhommeauParameters())
}

// computeHash computes a hash for the Delhommeau configuration
func (d *Delhommeau) computeHash() uint64 {
	h := fnv.New64a()
	// Sort keys for deterministic hash
	var keys []string
	for k := range d.exportableSettings {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(h, "%s=%v", k, d.exportableSettings[k])
	}
	return h.Sum64()
}

// Hash returns the hash of the Delhommeau configuration
func (d *Delhommeau) Hash() uint64 {
	return d.hash
}

// String returns a string representation showing only non-default values
func (d *Delhommeau) String() string {
	defaults := DefaultDelhommeauParameters()
	var nonDefaults []string

	if d.parameters.TabulationNr != defaults.TabulationNr {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_nr=%d", d.parameters.TabulationNr))
	}
	if d.parameters.TabulationRmax != defaults.TabulationRmax {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_rmax=%.1f", d.parameters.TabulationRmax))
	}
	if d.parameters.TabulationNz != defaults.TabulationNz {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_nz=%d", d.parameters.TabulationNz))
	}
	if d.parameters.TabulationZmin != defaults.TabulationZmin {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_zmin=%.1f", d.parameters.TabulationZmin))
	}
	if d.parameters.TabulationNbIntegrationPoints != defaults.TabulationNbIntegrationPoints {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_nb_integration_points=%d", d.parameters.TabulationNbIntegrationPoints))
	}
	if d.parameters.TabulationGridShape != defaults.TabulationGridShape {
		nonDefaults = append(nonDefaults, fmt.Sprintf("tabulation_grid_shape=%s", d.parameters.TabulationGridShape))
	}
	if d.parameters.FiniteDepthMethod != defaults.FiniteDepthMethod {
		nonDefaults = append(nonDefaults, fmt.Sprintf("finite_depth_method=%s", d.parameters.FiniteDepthMethod))
	}
	if d.parameters.FiniteDepthPronyDecompositionMethod != defaults.FiniteDepthPronyDecompositionMethod {
		nonDefaults = append(nonDefaults, fmt.Sprintf("finite_depth_prony_decomposition_method=%s", d.parameters.FiniteDepthPronyDecompositionMethod))
	}
	if d.parameters.FloatingPointPrecision != defaults.FloatingPointPrecision {
		nonDefaults = append(nonDefaults, fmt.Sprintf("floating_point_precision=%s", d.parameters.FloatingPointPrecision))
	}
	if d.parameters.GfSingularities != defaults.GfSingularities {
		nonDefaults = append(nonDefaults, fmt.Sprintf("gf_singularities=%s", d.parameters.GfSingularities))
	}

	if len(nonDefaults) == 0 {
		return "Delhommeau()"
	}
	return fmt.Sprintf("Delhommeau(%s)", fmt.Sprintf("%v", nonDefaults))
}

// createTabulation creates the tabulation for Green function computation
func (d *Delhommeau) createTabulation() error {
	// TODO: Implement tabulation creation
	// This would involve creating lookup tables for efficient Green function evaluation
	return nil
}

// createOrLoadTabulation creates or loads tabulation from cache
func (d *Delhommeau) createOrLoadTabulation() error {
	// TODO: Implement tabulation loading/creation with caching
	// Check if tabulation exists in cache directory, load if exists, create if not
	if d.parameters.TabulationCacheDir != "" {
		cacheFile := filepath.Join(d.parameters.TabulationCacheDir, fmt.Sprintf("tabulation_%d.cache", d.hash))
		_ = cacheFile // TODO: Implement actual caching logic
	}
	return d.createTabulation()
}

// Evaluate computes the Green function between two meshes using Delhommeau method
func (d *Delhommeau) Evaluate(mesh1, mesh2 interface{}, freeSurface float64, waterDepth float64,
	wavenumber complex128, adjointDoubleLayer bool, earlyDotProduct bool) (*mat.CDense, *mat.CDense, error) {

	// Get collocation points and normals
	colocationPoints, earlyDotProductNormals, err := d.getColocationPointsAndNormals(mesh1, mesh2, adjointDoubleLayer)
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
	S, K, err := d.initMatrices(rows, cols, earlyDotProduct)
	if err != nil {
		return nil, nil, err
	}

	// TODO: Implement actual Green function computation
	// This would involve:
	// 1. Computing distances between collocation points and mesh elements
	// 2. Evaluating Green function using tabulated values or direct computation
	// 3. Handling free surface effects and finite depth corrections
	// 4. Applying singularity treatments

	_ = colocationPoints
	_ = earlyDotProductNormals
	_ = freeSurface
	_ = waterDepth
	_ = wavenumber

	return S, K, nil
}

// GetParameters returns the current parameters
func (d *Delhommeau) GetParameters() DelhommeauParameters {
	return d.parameters
}
