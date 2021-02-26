/*
* @Author: 余添能
* @Date:   2021/2/26 8:57 上午
*/
package outer

type City struct {
	ID int
	Name string
	Spell string
}

type CityList struct {
	Initials string
	Cities []*City
}
