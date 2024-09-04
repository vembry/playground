package vembry.playground.app.models;

import java.math.BigDecimal;
import java.util.Date;

import com.fasterxml.jackson.annotation.JsonProperty;

import jakarta.persistence.Entity;
import jakarta.persistence.Id;
import jakarta.persistence.Table;

@Entity
@Table(name = "balances")
public class Balance {
    @Id
    @JsonProperty("id")
    String id;

    @JsonProperty("amount")
    BigDecimal amount;

    @JsonProperty("created_at")
    Date createdAt;

    @JsonProperty("updated_at")
    Date updatedAt;

    public Balance() {
    }

    public String getId() {
        return id;
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
