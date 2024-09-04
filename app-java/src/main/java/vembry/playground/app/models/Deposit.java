package vembry.playground.app.models;

import java.math.BigDecimal;
import java.time.Instant;
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

    @JsonProperty("status")
    @Enumerated(EnumType.STRING)
    Status status;

    @JsonProperty("amount")
    BigDecimal amount;

    @JsonProperty("created_at")
    Instant createdAt;

    @JsonProperty("updated_at")
    Instant updatedAt;

    public Deposit() {
    }

    public Deposit(String balanceId, BigDecimal amount) {
        this.id = Ksuid.newKsuid().toString();
        this.balanceId = balanceId;
        this.amount = amount;
        this.status = Status.pending;
        this.createdAt = Instant.now();
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

    public Instant getCreatedAt() {
        return createdAt;
    }

    public Instant getUpdatedAt() {
        return updatedAt;
    }
}
