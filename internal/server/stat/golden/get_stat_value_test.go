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
	"path"
	"runtime"
	"testing"

	pb "github.com/datacommonsorg/mixer/internal/proto"
	"github.com/datacommonsorg/mixer/test/e2e"
	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func TestGetStatValue(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	client, _, err := e2e.Setup()
	if err != nil {
		t.Fatalf("Failed to set up mixer and client")
	}
	_, filename, _, _ := runtime.Caller(0)
	goldenPath := path.Join(
		path.Dir(filename), "get_stat_value")

	for _, c := range []struct {
		statVar    string
		place      string
		goldenFile string
		mmethod    string
		wantErr    bool
	}{
		{
			"Count_Person",
			"country/USA",
			"count_person.json",
			"",
			false,
		},
		{
			"Count_CriminalActivities_CombinedCrime",
			"geoId/06",
			"total_crimes.json",
			"",
			false,
		},
		{
			"Annual_Generation_Electricity",
			"geoId/06",
			"electricity_generation.json",
			"",
			false,
		},
		{
			"Median_Age_Person",
			"geoId/0649670",
			"median_age.json",
			"",
			false,
		},
		{
			"Amount_EconomicActivity_GrossNationalIncome_PurchasingPowerParity_PerCapita",
			"country/USA",
			"gdp.json",
			"",
			false,
		},
		{
			"Count_Person",
			"country/USA",
			"empty.json",
			"bad_mmethod",
			true,
		},
		{
			"BadStatsVar",
			"geoId/06",
			"",
			"",
			true,
		},
		{
			"Count_Person",
			"badPlace",
			"",
			"",
			true,
		},
	} {
		resp, err := client.GetStatValue(ctx, &pb.GetStatValueRequest{
			StatVar:           c.statVar,
			Place:             c.place,
			MeasurementMethod: c.mmethod,
		})
		if c.wantErr {
			if err == nil {
				t.Errorf("Expect GetStatValue to error out but it succeed")
			}
			continue
		}
		if err != nil {
			t.Errorf("could not GetStatValue: %s", err)
			continue
		}
		if e2e.GenerateGolden {
			e2e.UpdateGolden(resp, goldenPath, c.goldenFile)
			continue
		}
		var expected pb.GetStatValueResponse
		if err = e2e.ReadJSON(goldenPath, c.goldenFile, &expected); err != nil {
			t.Errorf("Can not Unmarshal golden file")
			continue
		}

		if diff := cmp.Diff(resp, &expected, protocmp.Transform()); diff != "" {
			t.Errorf("payload got diff: %v", diff)
			continue
		}
	}
}
