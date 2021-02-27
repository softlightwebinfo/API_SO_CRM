package interfaces

type Method interface {
	Select()
	From()
	Join()
	Where()
	Order()
	Offset()
	Limit()
	Find()
	FindOne()
	FindAll()
}
