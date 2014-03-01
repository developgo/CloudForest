package CloudForest

import (
	"fmt"
)

/*
DensityTarget is used for density estimating trees. It contains a set of features and the
count of cases.
*/
type DensityTarget struct {
	Features *[]Feature
	N        int
}

/*
DensityTarget.SplitImpurity is a density estimating version of SplitImpurity.
*/
func (target *DensityTarget) SplitImpurity(l []int, r []int, m []int, counter *[]int) (impurityDecrease float64) {
	nl := float64(len(l))
	nr := float64(len(r))
	nm := 0.0

	impurityDecrease = nl * target.Impurity(&l, counter)
	impurityDecrease += nr * target.Impurity(&r, counter)
	if m != nil {
		nm := float64(len(m))
		impurityDecrease += nm * target.Impurity(&m, counter)
	}

	impurityDecrease /= nl + nr + nm
	return
}

//DensityTarget.Impurity uses the impurity measure defined in "Density Estimating Trees"
//by Parikshit Ram and Alexander G. Gray
func (target *DensityTarget) Impurity(cases *[]int, counter *[]int) (e float64) {
	t := len(*cases)
	e = float64(t*t) / float64(target.N*target.N)
	for _, f := range *target.Features {
		switch f.(type) {
		case CatFeature:
			bigenoughcounter := make([]int, f.NCats())
			e /= float64(f.(CatFeature).DistinctCats(cases, &bigenoughcounter))
		case NumFeature:
			e /= f.(NumFeature).Span(cases)
		}
	}

	return
}

//DensityTarget.FindPredicted returns the string representation of the density in the region
//spaned by the specified cases.
func (target *DensityTarget) FindPredicted(cases []int) string {
	t := len(cases)
	e := float64(t) / float64(target.N)

	for _, f := range *target.Features {
		switch f.(type) {
		case CatFeature:
			counter := make([]int, f.NCats())
			e /= float64(f.(CatFeature).DistinctCats(&cases, &counter))
		case NumFeature:
			e /= f.(NumFeature).Span(&cases)
		}
	}

	return fmt.Sprintf("%v", e)
}

func (target *DensityTarget) NCats() int {
	return 0
}