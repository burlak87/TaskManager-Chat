package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/your-team/taskmanager-chat/backend/internal/domain"
	"github.com/your-team/taskmanager-chat/backend/internal/service"
)

type BoardHandler struct {
	boardService *service.BoardService
}

func NewBoardHandler(boardService *service.BoardService) *BoardHandler {
	return &BoardHandler{
		boardService: boardService,
	}
}

func (h *BoardHandler) RegisterRoutes(router *gin.RouterGroup) {
	router.POST("/boards/:board_id/columns", h.createColumn)
	router.GET("/boards/:board_id/columns", h.getColumns)
	router.PUT("/boards/:board_id/columns/:column_id", h.updateColumn)
	router.DELETE("/boards/:board_id/columns/:column_id", h.deleteColumn)

	router.POST("/boards/:board_id/tasks", h.createTask)
	router.GET("/boards/:board_id/tasks", h.getTasks)
	router.GET("/boards/:board_id/tasks/:task_id", h.getTask)
	router.PUT("/boards/:board_id/tasks/:task_id", h.updateTask)
	router.DELETE("/boards/:board_id/tasks/:task_id", h.deleteTask)

	router.POST("/boards", h.createBoard)
	router.GET("/boards", h.getBoards)
	router.GET("/boards/:board_id", h.getBoard)
	router.PUT("/boards/:board_id", h.updateBoard)
	router.DELETE("/boards/:board_id", h.deleteBoard)
}

func (h *BoardHandler) createBoard(c *gin.Context) {
	var req domain.BoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userIDInt64 := userID.(int64)

	board, err := h.boardService.CreateBoard(c.Request.Context(), req.Title, req.Description, userIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, board)
}

func (h *BoardHandler) getBoards(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userIDInt64 := userID.(int64)

	boards, err := h.boardService.GetBoardsByOwner(c.Request.Context(), userIDInt64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, boards)
}

func (h *BoardHandler) getBoard(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	board, err := h.boardService.GetBoardByID(c.Request.Context(), boardID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "board not found"})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) updateBoard(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var req domain.BoardRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	board, err := h.boardService.UpdateBoard(c.Request.Context(), boardID, req.Title, req.Description)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, board)
}

func (h *BoardHandler) deleteBoard(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	err = h.boardService.DeleteBoard(c.Request.Context(), boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "board deleted"})
}

func (h *BoardHandler) createColumn(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var req domain.ColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	column, err := h.boardService.CreateColumn(c.Request.Context(), boardID, req.Title, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, column)
}

func (h *BoardHandler) getColumns(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	columns, err := h.boardService.GetColumnsByBoardID(c.Request.Context(), boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, columns)
}

func (h *BoardHandler) updateColumn(c *gin.Context) {
	columnIDStr := c.Param("column_id")
	columnID, err := strconv.ParseInt(columnIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	var req domain.ColumnRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	column, err := h.boardService.UpdateColumn(c.Request.Context(), columnID, req.Title, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, column)
}

func (h *BoardHandler) deleteColumn(c *gin.Context) {
	columnIDStr := c.Param("column_id")
	columnID, err := strconv.ParseInt(columnIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid column id"})
		return
	}

	err = h.boardService.DeleteColumn(c.Request.Context(), columnID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "column deleted"})
}

func (h *BoardHandler) createTask(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	var req domain.TaskRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.boardService.CreateTask(c.Request.Context(), boardID, req.ColumnID, req.Title, req.Description, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, task)
}

func (h *BoardHandler) getTasks(c *gin.Context) {
	boardIDStr := c.Param("board_id")
	boardID, err := strconv.ParseInt(boardIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid board id"})
		return
	}

	tasks, err := h.boardService.GetTasksByBoardID(c.Request.Context(), boardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *BoardHandler) getTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	task, err := h.boardService.GetTaskByID(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *BoardHandler) updateTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	var req domain.TaskUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := h.boardService.UpdateTask(c.Request.Context(), taskID, req.Title, req.Description, req.ColumnID, req.Position)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, task)
}

func (h *BoardHandler) deleteTask(c *gin.Context) {
	taskIDStr := c.Param("task_id")
	taskID, err := strconv.ParseInt(taskIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid task id"})
		return
	}

	err = h.boardService.DeleteTask(c.Request.Context(), taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "task deleted"})
}

