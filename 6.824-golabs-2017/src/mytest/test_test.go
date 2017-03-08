package mytest

import (
	"testing"
	"fmt"
)

type ByteSize float64

const (
	_ = iota
	KB ByteSize = 1 << (10 * iota)
	MB
	GB
	TB
)

func Test1(t *testing.T) {
	fmt.Printf("KB=%v, MB=%v, GB=%v, TB=%v", KB, MB, GB, TB)
}

func TestSlice(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice1Tmp := append(slice1, 4, 5)
	fmt.Printf("slice1==%v, slice1Tmp==%v\n", slice1, slice1Tmp)
	fmt.Printf("len slice1==%v, len slice1Tmp==%v\n", len(slice1), len(slice1Tmp))
	fmt.Printf("cap slice1==%v, cap slice1Tmp==%v\n", cap(slice1), cap(slice1Tmp))

	slice2 := make([]int, 10, 100)
	slice2Tmp := append(slice2, 4, 5)
	fmt.Printf("slice2==%v, slice2Tmp==%v\n", slice2, slice2Tmp)
	fmt.Printf("len slice2==%v, len slice2Tmp==%v\n", len(slice2), len(slice2Tmp))
	fmt.Printf("cap slice2==%v, cap slice2Tmp==%v\n", cap(slice2), cap(slice2Tmp))

	changeSlice(slice1)
	changeSlicePtr(&slice1)
}

func TestSlice1(t *testing.T) {
	slice1 := []int{1, 2, 3}
	slice1Tmp := append(slice1, 4, 5)
	fmt.Printf("slice1==%v, slice1Tmp==%v\n", slice1, slice1Tmp)
	fmt.Printf("len slice1==%v, len slice1Tmp==%v\n", len(slice1), len(slice1Tmp))
	fmt.Printf("cap slice1==%v, cap slice1Tmp==%v\n\n", cap(slice1), cap(slice1Tmp))

	slice2 := make([]int, 10, 100)
	slice2Tmp := append(slice2, 4, 5)
	fmt.Printf("slice2==%v, slice2Tmp==%v\n", slice2, slice2Tmp)
	fmt.Printf("len slice2==%v, len slice2Tmp==%v\n", len(slice2), len(slice2Tmp))
	fmt.Printf("cap slice2==%v, cap slice2Tmp==%v\n\n", cap(slice2), cap(slice2Tmp))

	/*changeSlice(slice1)
	changeSlicePtr(&slice1)*/
}

func TestSlicechange(t *testing.T) {
	slice1 := []int{1, 2, 3}
	fmt.Printf("slice1 Before==%v\n", slice1)
	changeSlice(slice1)
	fmt.Printf("slice1 After==%v\n\n", slice1)

	slice2 := []int{1, 2, 3}
	fmt.Printf("slice2 Before==%v\n", slice2)
	changeSlice(slice2)
	fmt.Printf("slice2 After==%v\n", slice2)
}

func TestSlicechange1(t *testing.T) {
	slice1 := make([]int, 10, 100)
	slice1[0] = 6666
	fmt.Printf("slice1 Before==%v slice1==%T slice1==%T\n  ", slice1, slice1, slice1[0])
	changeSlice(slice1)
	fmt.Printf("slice1 After==%v\n\n", slice1)

	slice2 := make([]int, 10, 100)
	slice2[0] = 6666
	fmt.Printf("slice2 Before==%v\n", slice2)
	changeSlice(slice2)
	fmt.Printf("slice2 After==%v\n", slice2)

}

func changeSlice(intSlice []int) {
	// 由于intslice是指向数组的指针
	// 因此操作指针映射的地址内容将会直接影响到原始数组
	intSlice[0] = 8888
}

func changeSlicePtr(intSlicePtr *[]int) {
	// 此处*intSlicePtr即为指向数组的指针
	slice := *intSlicePtr
	slice[0] = 9999
}

func TestMyAppend(t *testing.T) {
	s1 := make([]int, 10, 100)
	s1[0] = 1
	s1[1] = 2
	s1[2] = 3
	s1[3] = 4
	myAppend1(s1, []int{5, 6, 7})
	fmt.Printf("myAppend == %v\n", s1)

	s2 := make([]int, 10, 100)
	s2[0] = 1
	s2[1] = 2
	s2[2] = 3
	s2[3] = 4
	myAppend2(&s2, []int{5, 6, 7})
	fmt.Printf("myAppend1 == %v\n", s2)

	s3 := make([]int, 10, 100)
	s3[0] = 1
	s3[1] = 2
	s3[2] = 3
	s3[3] = 4
	myAppend3(&s3, []int{5, 6, 7})
	fmt.Printf("myAppend2 == %v\n", s3)
}

/**
 * 本次测试可知，数组定义形式为[3]int{},而切片为[]int{0,1,2}
 * 数组名不是指针，切片名为指向数组的指针
 * @author tantexian(https://my.oschina.net/tantexian/blog)
 * @since 2017/3/8
 * @params
 */
func TestArrayAndSlice(t *testing.T) {
	s0 := make([]int, 10, 100)
	s0[0] = 1
	s0[1] = 2
	s0[2] = 3
	s0[3] = 4
	passSlice(s0)
	fmt.Printf("slice ptr Before == %p\n\n", s0)

	a0 := [3]int{}
	a0[0] = 0
	a0[1] = 1
	a0[2] = 2
	passArray(a0)
	fmt.Printf("array a0  == %v\n", a0)
	fmt.Printf("array a0 type == %T\n", a0)
	fmt.Printf("array &a0 type == %T\n", &a0)
	fmt.Printf("array ptr Before == %p\n\n", &a0)

	a1 := []int{0, 1, 2}
	fmt.Printf("slice a1 == %v\n", a1)
	fmt.Printf("slice a1 type == %T\n", a1)
	fmt.Printf("slice &a1 type == %T\n", &a1)
	fmt.Printf("slice a1 ptr == %p\n", a1)
}

func passSlice(slice []int) {
	fmt.Printf("slice ptr After==%p \n", slice)
}

func passArray(array [3]int) {
	fmt.Printf("array ptr After==%p \n", &array)
}

func myAppend1(slice []int, data[]int) []int {
	// 此处的slice为指向数组的指针
	// 因此操作指针映射的地址内容将会直接影响到原始数组
	slice = []int{8, 8, 8, 8}
	slice[0] = 9999
	slice2 := data
	fmt.Printf("slice2==%v data=%v\n", slice2, data)
	fmt.Printf("slice2 ptr==%p data ptr==%p\n", slice2, data)
	return slice
}

func myAppend2(slicePtr *[]int, data[]int) []int {
	// 此处的slicePtr为指向数组指针的指针
	// *slicePtr为指向数组的指针
	// 赋值语句即将8,8,8,8的数组内容所在的地址赋值给*slicePtr，
	// 即*slicePtr指向8,8,8,8内容，因此slice原始内容改变了
	*slicePtr = []int{8, 8, 8, 8}
	(*slicePtr)[0] = 9999
	return *slicePtr
}

func myAppend3(slicePtr *[]int, data[]int) []int {
	// 此处slicePtr为指向指针的指针，该值为赋值指向8.8.8.8的内容
	// 原始的指针指向的内容仍然没有改变
	slicePtr = &[]int{8, 8, 8, 8}
	(*slicePtr)[0] = 9999
	return *slicePtr
}

func TestMap(t *testing.T) {
	map1 := make(map[int]int)
	map1[1] = 1
	changeMap1(map1)
	fmt.Printf("map1==%v\n", map1)
}

func changeMap1(intMap map[int]int) {
	intMap[1] = 2
}

type Counter int

func (ctr *Counter) ServeHTTP() {
	*ctr++
	fmt.Printf("counter = %d\n", *ctr)
}

func TestCount(t *testing.T) {
	count := new(Counter)
	count.ServeHTTP()

}