package handlers

import (
	"net/http"

	"wallet/internal/models"
	"wallet/internal/services"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	WalletService services.WalletService
}

type DepositRequest struct {
	Amount float64 `json:"amount"`
}
type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}
type TransferRequest struct {
	ToUserID string  `json:"to_user_id"`
	Amount   float64 `json:"amount"`
}

type TransactionHistoryRequest struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Type     string `form:"type"`
	Status   string `form:"status"`
}

type BalanceResponse struct {
	Balance float64 `json:"balance"`
}

type TransactionResponse struct {
	Balance float64 `json:"balance"`
}

type TransactionHistoryResponse struct {
	Transactions []models.Transaction `json:"transactions"`
}

func NewWalletHandler(walletService services.WalletService) *WalletHandler {
	return &WalletHandler{
		WalletService: walletService,
	}
}

func (h *WalletHandler) Deposit(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var req DepositRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.WalletService.Deposit(user.ID, req.Amount)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Balance: balance,
	})
}

func (h *WalletHandler) Withdraw(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var req WithdrawRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.WalletService.Withdraw(user.ID, req.Amount)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Balance: balance,
	})
}

func (h *WalletHandler) Transfer(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var req TransferRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	balance, err := h.WalletService.Transfer(user.ID, req.ToUserID, req.Amount)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Balance: balance,
	})
}

func (h *WalletHandler) GetBalance(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	balance, err := h.WalletService.GetBalance(user.ID)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, TransactionResponse{
		Balance: balance,
	})
}

func (h *WalletHandler) GetTransactionHistory(c *gin.Context) {
	user := c.MustGet("user").(*models.User)
	var req TransactionHistoryRequest

	if err := c.ShouldBindQuery(&req); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	transactions, err := h.WalletService.GetTransactionHistory(user.ID, req.Page, req.PageSize, req.Type, req.Status)
	if err != nil {
		c.AbortWithStatusJSON(err.Code, gin.H{"error": err.Message})
		return
	}

	c.JSON(http.StatusOK, TransactionHistoryResponse{Transactions: transactions})
}
