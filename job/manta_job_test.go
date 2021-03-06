package job_test

import (
	"testing"

	"github.com/image-server/image-server/job"
	. "github.com/image-server/image-server/test"
	"github.com/image-server/image-server/uploader/manta/client"
)

type FakeMantaClient struct {
	Opts client.CreateJobOpts
}

func (f *FakeMantaClient) CreateJob(opts client.CreateJobOpts) (string, error) {
	f.Opts = opts
	return "", nil
}

func TestMantaJobImageCommand(t *testing.T) {
	job := job.MantaJob{BasePath: "tacos", Outputs: "whatever", Namespace: "p"}
	var output = job.ToImageCommand()
	Matches(t, "/assets/wanelo/public/images/bin/images-solaris-1", output)
	Matches(t, "--remote_base_path tacos --namespace p --outputs whatever process", output)
	Matches(t, "MANTA_INPUT_FILE", output)
}

func TestMantaJobOpts(t *testing.T) {
	job := job.MantaJob{BasePath: "tacos", Outputs: "whatever"}
	var opts = job.ToJobOpts()
	Matches(t, "bin/images-solaris-1.", opts.Phases[0].Exec)
	Equals(t, client.Phase{Type: "reduce", Exec: "cat"}, opts.Phases[1])
}

func TestMantaJobCreation(t *testing.T) {
	job := job.MantaJob{BasePath: "tacos", Outputs: "whatever"}
	fakeMantaClient := &FakeMantaClient{}
	_, err := job.CreateMantaJob(fakeMantaClient)
	Ok(t, err)
	Equals(t, job.ToJobOpts(), fakeMantaClient.Opts)
}
