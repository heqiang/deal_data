package datadeal

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"
)

type AutoGenerated []struct {
	Type        string      `json:"type"`
	Name        interface{} `json:"name,omitempty"`
	Md5Src      string      `json:"md5src,omitempty"`
	Description string      `json:"description,omitempty"`
	Src         string      `json:"src,omitempty"`
	Data        string      `json:"data,omitempty"`
}

func TestNewDataDeal(t *testing.T) {
	testContents := `[{"type": "image", "name": null, "md5src": "f4bee9790cce0145e2bd571b12ec78b3.jpg", "description": "Music video director Shin Hee-won has apologised for plagiarising Tokyo DisneySea's 15th anniversary design.", "src": "https://media.asiaone.com/sites/default/files/styles/article_main_image/public/original_images/Aug2022/150822_screengrabforever1_smentdisney2%20%281%29%20%281%29.jpeg?itok=DhKuqNb0"}, {"data": "Stay in the know with a recap of our top stories today.", "type": "text"}, {"data": "1.", "type": "text"}, {"data": "'Exactly the same': Girls' Generation's new music video struck by plagiarism claims, director apologises", "type": "text"}, {"data": "Their new music video is supposed to celebrate their return but this mistake might have left them flying too close to the sun... \u00bb", "type": "text"}, {"data": "READ MORE", "type": "text"}, {"data": "2.", "type": "text"}, {"data": "'The boy admitted and apologised': Alleged thief returns empty purse which contained $1,690 after stealing from single mum of 5 in Punggol", "type": "text"}, {"type": "image", "name": null, "md5src": "5ad4f3913bfcb00069f9379d2a214aa6.jpg", "description": null, "src": "https://media.asiaone.com/sites/default/files/styles/article_main_image/public/original_images/Aug2022/170822_stomp_stomp.jpg?itok=rkumOwms"}, {"data": "Stomp contributor Nora (not her real name) earlier shared how her purse had been stolen by a teenager at a playground... \u00bb", "type": "text"}, {"data": "READ MORE", "type": "text"}, {"data": "3.", "type": "text"}, {"data": "Elderly artist at Toa Payoh tears after making her first sale in days", "type": "text"}, {"type": "image", "name": null, "md5src": "61d1b3345d2529a6be3f5b927c62806b.jpg", "description": null, "src": "https://media.asiaone.com/sites/default/files/styles/article_main_image/public/original_images/Aug2022/160822_crotchet_HLfb.jpg?itok=9fD0Lr_e"}, {"data": "Jaya Dutta was at Toa Payoh MRT station last Friday (Aug 12) morning and came across an elderly woman whom she sees on the way to work...\u00a0\u00bb", "type": "text"}, {"data": "READ MORE", "type": "text"}, {"data": "4.", "type": "text"}, {"data": "'Lesson learnt, never polish your cars': Peacock attacks its own reflection on vehicle in Sentosa", "type": "text"}, {"type": "image", "name": null, "md5src": "bf9962e688571e26be49df712061078a.jpg", "description": null, "src": "https://media.asiaone.com/sites/default/files/styles/article_main_image/public/original_images/Aug2022/peacock.jpg?itok=MCo7tJpD"}, {"data": "A shiny and reflective navy blue car was seen getting attacked by a peacock in Sentosa... \u00bb", "type": "text"}, {"data": "READ MORE", "type": "text"}, {"data": "editor@asiaone.com", "type": "text"}]`

	d := NewDataDeal("", "cn")
	newContent := d.transContent(testContents)
	fmt.Println(newContent)
}

func TestFileNmae(t *testing.T) {
	filePath := `E:\\goproject\\deal_data\\data\\2022-09-20\\20220920202202_newsty.json`
	res2 := strings.Split(filePath, string(filepath.Separator))
	filep := strings.Join(res2[0:len(res2)-1], string(filepath.Separator))
	println(filep)
}
