package main

// func main() {
// 	zipPath := "games.zip" // the path to the .zip folder

// 	r, err := zip.OpenReader(zipPath)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	defer r.Close()

// 	for _, f := range r.File {
// 		if !f.FileInfo().IsDir() && filepath.Ext(f.Name) == ".mp4" {
// 			fmt.Println(f.Name)
// 		}
// 	}
// }
