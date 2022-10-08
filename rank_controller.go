package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"time"
)

const RANK_KEY = "CurrentMonthRank"

type RankInfo struct {
	score int64
	id    int
	rank  int
}

func QueryTopThree(client redis.Conn) []RankInfo {
	return QueryRankInfo(client, 0, 3)
}

func QueryAroundTenRank(client redis.Conn, id int) []RankInfo {
	myrank, err := redis.Int(client.Do("zrevrank", RANK_KEY, id))
	if err != nil {
		fmt.Println("zincrby failed, err:", err)
		return nil
	}
	var beginRank, endRank int
	//排名在前5，取前11个，否则取前后各5个
	if myrank < 5 {
		beginRank = 0
		endRank = 10
	} else {
		beginRank = myrank - 5
		endRank = myrank + 6
	}
	rankInfos := QueryRankInfo(client, beginRank, endRank)
	for i := 0; i < len(rankInfos); i++ {
		if rankInfos[i].rank == myrank {
			rankInfos = append(rankInfos[:i], rankInfos[i+1:]...)
			continue
		}
	}
	return rankInfos
}

//查询排名信息，相同分数合成一位
func QueryRankInfo(client redis.Conn, beginRank, endRank int) []RankInfo {
	var rankInfos []RankInfo
	res, err := redis.Ints(client.Do("zrevrange", RANK_KEY, beginRank, endRank))
	if err != nil {
		fmt.Println("zincrby failed, err:", err)
		return nil
	}
	rank := beginRank
	for i := 0; i < len(res); i++ {
		rank++
		rankInfo := RankInfo{
			rank: rank,
			id:   res[i],
		}
		rankInfos = append(rankInfos, rankInfo)
	}
	for i := 0; i < len(rankInfos); i++ {
		score, err := redis.Int(client.Do("zscore", RANK_KEY, rankInfos[i].id))
		if err != nil {
			fmt.Println("zincrby failed, err:", err)
			return nil
		}
		rankInfos[i].score = getPoint(int64(score))
	}
	nowRank := rankInfos[0].rank
	for i := 1; i < len(rankInfos); i++ {
		if rankInfos[i].score != rankInfos[i-1].score {
			nowRank++
		}
		rankInfos[i].rank = nowRank
	}
	return rankInfos
}

//计算当月排行榜score值，用redis zset实现，分值高22位表示积分，第41位表示时间戳(结束时间戳-当前时间戳)
func toScore(point int64) int64 {
	var score int64
	score = 0
	periodEndTimestamp := getMonthLastDayMill()
	score = (score | point) << 41
	score = score | (periodEndTimestamp - time.Now().Unix())
	score = score | periodEndTimestamp
	return score
}

//获取实际积分
func getPoint(score int64) int64 {
	return score >> 41
}

//更新积分
func updateRank(client redis.Conn, id int64, point int64) {
	score := toScore(point)
	_, err := client.Do("zadd", RANK_KEY, score, id)
	if err != nil {
		fmt.Println("zadd failed, err:", err)
		return
	}
}

//查询排名
func queryMyRank(client redis.Conn, id int) int {
	res, err := redis.Int(client.Do("zrevrank", RANK_KEY, id))
	if err != nil {
		fmt.Println("zincrby failed, err:", err)
		return 0
	}
	return res
}

func main() {
	client, err := redis.Dial("tcp", "localhost:6379")
	defer client.Close()
	if err != nil {
		fmt.Println("redis connect failed,", err)
		return
	}
	fmt.Println("redis connect success")
	updateRank(client, 3, 999)
	updateRank(client, 1, 9)
	updateRank(client, 2, 99)
	time.Sleep(1e9)
	updateRank(client, 4, 999)
	res := QueryTopThree(client)
	fmt.Printf("top 3 is %v\n", res)
	ranks := QueryAroundTenRank(client, queryMyRank(client, 2))
	for i := 0; i < len(ranks); i++ {
		fmt.Printf("rankinfos %v\n", ranks[i])
	}
}

//获取当月最后一天time-mill
func getMonthLastDayMill() int64 {
	now := time.Now()
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	return time.Date(currentYear, currentMonth+1, -1, 0, 0, 0, 0, currentLocation).UnixMilli()
}
