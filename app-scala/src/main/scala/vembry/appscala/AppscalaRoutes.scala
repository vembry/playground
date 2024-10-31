package vembry.appscala

import cats.effect.Sync
import cats.syntax.all.*
import org.http4s.HttpRoutes
import org.http4s.dsl.Http4sDsl

object AppscalaRoutes:

  // Define routes for the joke service, which responds with a joke in JSON format
  def jokeRoutes[F[_]: Sync](J: Jokes[F]): HttpRoutes[F] =
    val dsl = new Http4sDsl[F]{}  // Provides HTTP DSL to handle requests and responses
    import dsl.*
    HttpRoutes.of[F] {  // Define the route as a partial function
      case GET -> Root / "joke" =>  // Matches GET requests to /joke
        for {
          joke <- J.get              // Retrieve a joke from the Jokes service
          resp <- Ok(joke)           // Respond with the joke in an HTTP 200 OK response
        } yield resp
    }

  // Define routes for the hello world service, which responds with a greeting in JSON format
  def helloWorldRoutes[F[_]: Sync](H: HelloWorld[F]): HttpRoutes[F] =
    val dsl = new Http4sDsl[F]{}
    import dsl.*
    HttpRoutes.of[F] {
      case GET -> Root / "hello" / name =>  // Matches GET requests to /hello/<name>
        for {
          greeting <- H.hello(HelloWorld.Name(name))  // Generate a greeting message
          resp <- Ok(greeting)                         // Respond with the greeting in JSON
        } yield resp
    }
