package cli



import "testing"

func TestCli(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping test in short mode.")
    }
}