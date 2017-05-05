package gocaptcha

type Point struct {
	X int
	Y int
}

func NewPoint(x int,y int) *Point{
	return &Point{ X :x,Y:y};
}
