package workflow

type Recorder struct{ Messages []string }

func (r *Recorder) Notify(message string) error { r.Messages = append(r.Messages, message); return nil }
func (r *Recorder) Last() string {
	if len(r.Messages) == 0 {
		return ""
	}
	return r.Messages[len(r.Messages)-1]
}
