package mediator

import "testing"

func TestMediator(t *testing.T) {
	mediator := GetMediatorInstance()
	mediator.CD = &CdDriver{}
	mediator.CPU = &Cpu{}
	mediator.Sound = &SoundCard{}
	mediator.Video = &VideoCard{}

	mediator.CD.ReadData()

	if mediator.CD.Data != "music,image" {
		t.Fatalf("CD unexpect data %s", mediator.CD.Data)
	}

	if mediator.CPU.Sound != "music" {
		t.Fatalf("CPU unexpect sound data %s", mediator.CPU.Sound)
	}

	//if mediator.CPU.Video != "image" {
	//	t.Fatalf("CPU unexpect video data %s", mediator.CPU.Video)
	//}

	//if mediator.Video.Data != "image" {
	//	t.Fatalf("VideoCard unexpect data %s", mediator.Video.Data)
	//}

	if mediator.Sound.Data != "music" {
		t.Fatalf("SoundCard unexpect data %s", mediator.Sound.Data)
	}

}
