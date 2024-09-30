package vembry.playground.app.models;

import java.math.BigDecimal;
import java.time.Instant;

import com.fasterxml.jackson.annotation.JsonProperty;
import com.github.ksuid.Ksuid;

import jakarta.persistence.Entity;
import jakarta.persistence.EnumType;
import jakarta.persistence.Enumerated;
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
    @Enumerated(EnumType.STRING)
    Status status;

    @JsonProperty("amount")
    BigDecimal amount;

    @JsonProperty("created_at")
    Instant createdAt;

    @JsonProperty("updated_at")
    Instant updatedAt;

    public Withdrawal() {
    }

    public Withdrawal(String balanceId, BigDecimal amount) {
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
