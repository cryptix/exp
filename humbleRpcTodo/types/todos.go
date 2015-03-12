package types

type Todo struct {
	Id          int
	Title       string
	IsCompleted bool
}

type TodosIndex map[int]*Todo

func (t *Todo) CheckedStr() string {
	if t.IsCompleted {
		return "checked"
	}
	return ""
}

func (t *Todo) CompletedStr() string {
	if t.IsCompleted {
		return "completed"
	}
	return ""
}

type TodoListArgs struct {
	What string
}

type ReplyOK struct {
	Id int
	OK bool
}

type ById []Todo

func (l ById) Len() int {
	return len(l)
}

func (l ById) Less(i int, j int) bool {
	return l[i].Id < l[j].Id
}

func (l ById) Swap(i int, j int) {
	l[i], l[j] = l[j], l[i]
}
