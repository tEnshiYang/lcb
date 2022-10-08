package main

func is_full(board [][]int, row int) bool {
	lenth := len(board[0])
	for i := 0; i < lenth; i++ {
		if board[row][i] == 0 {
			return false
		}
	}
	return true
}

func lock_dot(board [][]int, row int, col int) bool {
	//越界直接返回
	if row < 0 || row > len(board)-1 || col < 0 || col > len(board[0])-1 {
		return false
	}
	if board[row][col] == 1 {
		return false
	}
	return true
}

//row和col是鼠标左键松开落点位置
func lock_piece(board [][]int, row int, col int, piece string) bool {
	//获取指定图形所有坐标
	piecePos := get_piece_pos(row, col, piece)
	//判断图形每个坐标点是否可放置
	for i := 0; i < len(piecePos); i++ {
		if !lock_dot(board, piecePos[i][0], piecePos[i][1]) {
			return false
		}
	}
	return true
}

func get_piece_pos(row int, col int, piece string) [][]int {
	var res [][]int
	switch piece {
	case "Green_L":
		res = append(res, []int{row - 1, col + 1})
		res = append(res, []int{row - 1, col})
		res = append(res, []int{row - 1, col - 1})
		res = append(res, []int{row, col - 1})
		res = append(res, []int{row + 1, col - 1})
		return res
	default:
		return res
	}
}
