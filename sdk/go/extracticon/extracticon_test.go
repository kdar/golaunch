package extracticon

import (
	"fmt"
	"image/png"
	"os"
	"testing"

	"github.com/ansel1/merry"
)

func TestFrom(t *testing.T) {
	// images, err := From("D:\\Program Files (x86)\\SMPlayer\\smplayer.exe", []int{32})
	// if err != nil {
	// 	t.Fatal(merry.Details(err))
	// }
	//
	// fmt.Println(images)

	extract := New()

	image, err := extract.From("testdata/smush.exe")
	if err != nil || image == nil {
		t.Fatalf("%s: %v", "", merry.Details(err))
	}

	fp, _ := os.Create("output.png")
	defer fp.Close()

	err = png.Encode(fp, image)
	if err != nil {
		fmt.Println(err)
	}
}
