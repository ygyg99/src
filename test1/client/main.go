package main //包，表明代码所在模块
import (
	"fmt"
)

func find(n []int, target, idx int) int{
    l,r := 0, len(n)
    mid := (l+r+1)/2
    if mid >0 && n[mid-1] > n[mid]{
        return mid+idx
    }
    if mid == r{
        return 0
    }
    return find(n[l:mid],target,idx) + find(n[mid:r],target,idx+mid)
}
func search(nums []int, target int) int {
   return find(nums,target,0)
}

func main() {
	a := []int{4,5,6,7,0,1,2}
	fmt.Println(search(a,0))
	// 现在 p1.X 和 p1.Y 是 2 和 4
}
