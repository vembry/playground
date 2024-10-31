package vembry.appscala

import cats.effect.Async
import cats.syntax.all.*
import com.comcast.ip4s.*
import fs2.io.net.Network
import org.http4s.ember.client.EmberClientBuilder
import org.http4s.ember.server.EmberServerBuilder
import org.http4s.implicits.*
import org.http4s.server.middleware.Logger

object AppscalaServer:

  // Define a function to run the server, parameterized with an `F[_]` effect type
  def run[F[_]: Async: Network]: F[Nothing] = {
    for {
      // Build an HTTP client that will be used by the joke service to make external requests
      client <- EmberClientBuilder.default[F].build

      // Instantiate service implementations
      helloWorldAlg = HelloWorld.impl[F]
      jokeAlg = Jokes.impl[F](client)

      // Combine service routes into an HTTP application
      httpApp = (
        AppscalaRoutes.helloWorldRoutes[F](helloWorldAlg) <+>
        AppscalaRoutes.jokeRoutes[F](jokeAlg)
      ).orNotFound

      // Add logging middleware to log each request and response
      finalHttpApp = Logger.httpApp(true, true)(httpApp)

      // Start the server using Ember, binding to 0.0.0.0:8080
      _ <- 
        EmberServerBuilder.default[F]
          .withHost(ipv4"0.0.0.0")
          .withPort(port"8080")
          .withHttpApp(finalHttpApp)
          .build
    } yield ()
  }.useForever  // Keeps the server running indefinitely
