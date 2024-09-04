package vembry.playground.app.models;

import java.math.BigDecimal;
import java.util.Date;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.ksuid.Ksuid;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(name = "withdrawals")
public class Withdrawal {
    @Id
    @JsonProperty("id")
    String id;

    @JsonProperty("balance_id")
    String balanceId;

    @JsonProperty("status")
    String status;

    @JsonProperty("amount")
    BigDecimal amount;

    @JsonProperty("created_at")
    Date createdAt;

    @JsonProperty("updated_at")
    Date updatedAt;

    public Withdrawal() {
    }

    public Withdrawal(String balanceId, BigDecimal amount) {
        this.id = Ksuid.newKsuid().toString();
        this.balanceId = balanceId;
        this.amount = amount;
        this.status = "pending";
    }

    public String getId() {
        return id;
    }

    public String getBalanceId() {
        return balanceId;
    }

    public BigDecimal getAmount() {
        return amount;
    }

    public Date getCreatedAt() {
        return createdAt;
    }

    public Date getUpdatedAt() {
        return updatedAt;
    }
}
