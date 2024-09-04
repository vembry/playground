package vembry.playground.app.models;

import java.math.BigDecimal;

import com.fasterxml.jackson.annotation.JsonProperty;

public class WithdrawRequest {
    @JsonProperty("amount")
    BigDecimal amount;
    @JsonProperty("balance_id")
    String balanceId;

    public void setAmount(BigDecimal amount) {
        this.amount = amount;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public void setBalanceId(String balanceId) {
        this.balanceId = balanceId;
    }

    public String getBalanceId() {
        return balanceId;
    }
}
