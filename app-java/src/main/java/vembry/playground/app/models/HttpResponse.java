package vembry.playground.app.models;

import com.fasterxml.jackson.annotation.JsonProperty;

public class HttpResponse<T> {
    @JsonProperty("error")
    String error;
    @JsonProperty("data")
    T data;

    public HttpResponse(String error, T data) {
        this.error = error;
        this.data = data;
    }

    public HttpResponse(String error) {
        this.error = error;
    }
    
    public HttpResponse(T data) {
        this.data = data;
    }

    public String getError() {
        return error;
    }

    public T getData() {
        return data;
    }
}
