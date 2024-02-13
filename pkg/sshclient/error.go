package sshclient

type NotFoundError struct{}

func (o NotFoundError) Error() string {
	return "ONT not found"
}

type InvalidSerialNumberError struct{}

func (i InvalidSerialNumberError) Error() string {
	return "Invalid serial number"
}
