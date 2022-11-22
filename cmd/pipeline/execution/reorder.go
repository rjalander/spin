// Copyright (c) 2022, Google, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.

package execution

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/spf13/cobra"

)

var (
	reorderExecutionShort = "Re-order the waiting execution to UP/DOWN for the provided execution id"
	reorderExecutionLong  = "Re-order the waiting execution to UP/DOWN for the provided execution id"
)

type reorderOptions struct {
	*executionOptions
        output            string
	executionId       string
        reorderAction     string
}

func NewReorderCmd(executionOptions *executionOptions) *cobra.Command {
	options := &reorderOptions{
		executionOptions: executionOptions,
	}
	cmd := &cobra.Command{
		Use:   "reorder",
		Short: reorderExecutionShort,
		Long:  reorderExecutionLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return reorderExecution(cmd, options)
		},
	}
	cmd.PersistentFlags().StringVarP(&options.executionId, "execution-id","i", "", "Spinnaker waiting execution id to reorder")
        cmd.PersistentFlags().StringVarP(&options.reorderAction, "reorder-action","r", "", "Re-order to UP/DOWN for the execution id")

	return cmd
}

func reorderExecution(cmd *cobra.Command, options *reorderOptions) error {
	if options.executionId == "" {
                return errors.New("required parameter 'execution-id' not set")
        }
	if options.reorderAction == "" {
                return errors.New("required parameter 'reorder-action' not set")
        }
	if (!strings.EqualFold(options.reorderAction, "UP") && !strings.EqualFold(options.reorderAction, "Down")) {
                return errors.New("required parameter 'reorder-action' not set to UP/DOWN")
        }

	resp, err := options.GateClient.PipelineControllerApi.ReorderPipelineUsingPUT(options.GateClient.Context, options.executionId, options.reorderAction)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("encountered an error re-ordering execution with id %s, reorder-action %s, status code: %d\n",
			options.executionId,
			options.reorderAction,
			resp.StatusCode)
	}
	if err != nil {
		return err
	}

	options.Ui.Success(fmt.Sprintf("Execution %s successfully re-ordered to %s", options.executionId, options.reorderAction))
	return nil
}
