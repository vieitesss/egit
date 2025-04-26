package msgs

type GitCmdOutMsgI interface {
	FromGitCmdOut(GitCmdOutMsg) GitCmdOutMsgI
}

type GitCmdOutMsg struct {
	Out string
	Err error
}
