// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/google/pullsheet/pkg/client"
	"github.com/google/pullsheet/pkg/server/job"
)

type Server struct {
	cl   *client.Client
	jobs []*job.Job
}

func New(ctx context.Context, c *client.Client, initJob *job.Job) *Server {
	jobs := []*job.Job{}
	if initJob != nil {
		jobs = append(jobs, initJob)
		go initJob.Update(ctx, c)
	}

	return &Server{
		cl:   c,
		jobs: jobs,
	}
}

func (s *Server) Root() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, err := s.jobs[0].Render()
		if err != nil {
			logrus.Errorf("rendering job page: %s", err)
		}
		fmt.Fprint(w, res)
	}
}
