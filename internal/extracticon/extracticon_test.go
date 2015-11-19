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

	//image, err := extract.From("D:\\Program Files\\Adobe\\Adobe After Effects CC 2014\\Support Files\\AfterFX.exe")
	//image := From2("D:\\Program Files\\Microsoft Network Monitor 3\\NmBuild.exe")
	//image := From2("D:\\Program Files (x86)\\Elaborate Bytes\\CloneDVD2\\HelpLauncher.exe")
	image, err := extract.From("C:\\Users\\outroot\\AppData\\Roaming\\Microsoft\\Windows\\Start Menu\\RaidCall.lnk")
	// if err != nil || image == nil {
	// 	image, err = extract.FromExt(filepath.Ext("D:\\Program Files (x86)\\Zeus\\zGNU\\tidy.exe"))
	// }
	//image, err := extract.From("testdata/Icon300.ico")
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
