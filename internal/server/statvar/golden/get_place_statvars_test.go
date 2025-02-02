// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package golden

import (
	"context"
	"testing"

	pb "github.com/datacommonsorg/mixer/internal/proto"
	"github.com/datacommonsorg/mixer/test/e2e"
)

func TestGetPlaceStatVars(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, _, err := e2e.Setup(&e2e.TestOption{UseMemdb: true})
	if err != nil {
		t.Fatalf("Failed to set up mixer and client")
	}

	for _, c := range []struct {
		dcids    []string
		want     []string
		minCount int
		wanterr  bool
	}{
		{
			[]string{"geoId/05"},
			[]string{"Count_Person"},
			1000,
			false,
		},
		{
			[]string{"geoId/06085"},
			[]string{"Count_Person"},
			1000,
			false,
		},
		{
			[]string{"country/ALB"},
			[]string{"GiniIndex_EconomicActivity", "Test_Stat_Var_1"},
			0,
			false,
		},
		{
			[]string{"invalid"},
			[]string{},
			0,
			false,
		},
		{
			[]string{},
			[]string{},
			0,
			true,
		},
	} {

		req := &pb.GetPlaceStatVarsRequest{
			Dcids: c.dcids,
		}
		resp, err := client.GetPlaceStatVars(ctx, req)
		if c.wanterr {
			if err == nil {
				t.Errorf("Expect to get error for GetPlaceStatsVar() but succeed")
			}
			continue
		}
		if err != nil {
			t.Errorf("Could not GetPlaceStatsVar: %s", err)
			continue
		}
		for dcid, place := range resp.Places {
			if len(place.StatVars) < c.minCount {
				t.Errorf("%s has less than %d stats vars", dcid, c.minCount)
			}
			statsVarSet := map[string]bool{}
			for _, statsVar := range place.StatVars {
				statsVarSet[statsVar] = true
			}
			for _, statsVar := range c.want {
				if _, ok := statsVarSet[statsVar]; !ok {
					t.Errorf("%s is not in the stats var list of %s", statsVar, dcid)
					continue
				}
			}
		}
	}
}
