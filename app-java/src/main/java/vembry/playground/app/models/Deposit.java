package vembry.playground.app.models;

import java.math.BigDecimal;
import java.util.Date;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.ksuid.Ksuid;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;
import jakarta.persistence.Enumerated;
import jakarta.persistence.EnumType;

@Entity
@Table(name = "deposits")
public class Deposit {
    @Id
    @JsonProperty("id")
    String id;

    @JsonProperty("balance_id")
    String balanceId;

    @Enumerated(EnumType.STRING)
    @JsonProperty("status")
    Status status;

    @JsonProperty("amount")
    BigDecimal amount;

    @JsonProperty("created_at")
    Date createdAt;

    @JsonProperty("updated_at")
    Date updatedAt;

    public Deposit() {
    }

    public Deposit(String balanceId, BigDecimal amount) {
        this.id = Ksuid.newKsuid().toString();
        this.balanceId = balanceId;
        this.amount = amount;
        this.status = Status.pending;
    }

    public String getId() {
        return id;
    }

    public String getBalanceId() {
        return balanceId;
    }

    public Status getStatus() {
        return status;
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
