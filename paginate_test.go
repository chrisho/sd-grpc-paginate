package paginate

import (
	"fmt"
	"testing"
	"github.com/jinzhu/gorm"
)

//func TestGetPagingOptions(t *testing.T){
//	in := &PageOptions{
//		PageSize:-1,
//		PageNumber:4,
//		//CurrPageNumber:1,
//		//SortValue:2,
//	}
//	offset, limit,symbol := GetPagingOptions(in)
//	fmt.Println(offset,limit,symbol)
//}
func TestGetPaginate(t *testing.T){
	in := &PageOptions{
		PageSize:-1,
		PageNumber:4,
		//CurrPageNumber:1,
		//SortValue:2,
		SortField:"menu_id_a",
	}
	db := &gorm.DB{}
	sql := GetPaginate(db,"xy_role_menu",in)
	fmt.Println(sql.Debug())
}