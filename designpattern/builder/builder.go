package builder

import "fmt"

type Builder interface {
	EatFood()
	WearClothes()
	Sleep()
}

type Director struct {
	builder Builder
}

// ConstructA 构造器A
func (b *Director) ConstructA() {
	b.builder.EatFood()
	b.builder.Sleep()
}

// ConstructB 构造器B
func (b *Director) ConstructB() {
	b.builder.EatFood()
	b.builder.WearClothes()
	b.builder.Sleep()
}

func NewDirector(builder Builder) *Director {
	return &Director{
		builder: builder,
	}
}

type builderA struct {
	Food     string
	Clothes  string
	SleepBed string
}

func (b *builderA) EatFood() {
	b.Food = "kfc"
}

func (b *builderA) WearClothes() {
	b.Clothes = "nike"
}

func (b *builderA) Sleep() {
	b.SleepBed = "bed"
}

func (b builderA) ActionResult() string {
	return fmt.Sprintf("tom wear %s, eat %s and go to %s", b.Clothes, b.Food, b.SleepBed)
}
