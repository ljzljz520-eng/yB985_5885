package workflow

import "errors"

type Step struct {
	Code, Label string
	Done        bool
}

func OpenSteps() []Step {
	return []Step{{"open", "创建巡店记录", false}, {"finding", "登记问题", false}, {"assign", "分配负责人", false}, {"notify", "发送通知", false}}
}
func RemediationSteps() []Step {
	return []Step{{"write", "填写整改说明", false}, {"attach", "上传附件", false}, {"review", "审核整改", false}, {"close", "关闭记录", false}}
}
func Complete(steps []Step) bool {
	for _, s := range steps {
		if !s.Done {
			return false
		}
	}
	return true
}
func Mark(steps []Step, code string) error {
	for i := range steps {
		if steps[i].Code == code {
			steps[i].Done = true
			return nil
		}
	}
	return errors.New("unknown step")
}
func CanAdvance(current, next string) bool {
	seq := []string{"open", "finding", "assign", "notify"}
	for i, v := range seq {
		if v == current && i+1 < len(seq) && seq[i+1] == next {
			return true
		}
	}
	return false
}
