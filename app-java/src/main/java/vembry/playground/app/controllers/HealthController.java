package vembry.playground.app.controllers;

import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import vembry.playground.app.models.HttpResponse;

@RestController
@RequestMapping("/health")
public class HealthController {

    @GetMapping
    public HttpResponse<Object> get() {
        return new HttpResponse<Object>(null);
    }
}
