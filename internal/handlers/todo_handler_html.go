package handlers

import (
	"gotta/internal/httpx"
	"gotta/internal/services"
	"gotta/web/components"
	"gotta/web/viewmodels"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TodoHandlerHTML struct {
	todoService *services.TodoService
}

func NewTodoHandlerHTML(todoService *services.TodoService) *TodoHandlerHTML {
	return &TodoHandlerHTML{todoService: todoService}
}

func (h *TodoHandlerHTML) Index(c *gin.Context) {
	todo, err := h.todoService.TodoList()

	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Failed to load todos")
		return
	}

	meta := viewmodels.PageMeta{
		Title:       "Todo List",
		Description: "A simple todo list app",
	}

	httpx.Render(c, http.StatusOK, components.TodoPage(todo, meta))
}

func (h *TodoHandlerHTML) Add(c *gin.Context) {
	text := c.PostForm("text")

	todos, err := h.todoService.Add(text)
	if err != nil {
		c.String(http.StatusBadRequest, err.Error())
		return
	}

	httpx.Render(c, http.StatusOK, components.TodoList(todos))
}

func (h *TodoHandlerHTML) Toggle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todos, err := h.todoService.Toggle(id)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Failed to toggle todo")
		return
	}

	httpx.Render(c, http.StatusOK, components.TodoList(todos))

}

func (h *TodoHandlerHTML) Delete(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.String(http.StatusBadRequest, "Invalid todo ID")
		return
	}

	todos, err := h.todoService.Delete(id)
	if err != nil {
		c.Error(err)
		c.String(http.StatusInternalServerError, "Failed to delete todo")
		return
	}

	httpx.Render(c, http.StatusOK, components.TodoList(todos))
}
