package vembry.appscala

import cats.Applicative
import cats.syntax.all.*
import io.circe.{Encoder, Json}
import org.http4s.EntityEncoder
import org.http4s.circe.*

trait HelloWorld[F[_]]:
  // Defines a method to produce a greeting for a given name
  def hello(n: HelloWorld.Name): F[HelloWorld.Greeting]

object HelloWorld:
  // Wrapper case class to represent a name
  final case class Name(name: String) extends AnyVal

  // Wrapper case class to represent a greeting message
  final case class Greeting(greeting: String) extends AnyVal

  object Greeting:
    // Encoder for converting Greeting into JSON
    given Encoder[Greeting] = new Encoder[Greeting]:
      final def apply(a: Greeting): Json = Json.obj(
        ("message", Json.fromString(a.greeting)),
      )

    // Entity encoder for converting Greeting to an HTTP response
    given [F[_]]: EntityEncoder[F, Greeting] = jsonEncoderOf[F, Greeting]

  // Implementation of the HelloWorld service that provides greeting messages
  def impl[F[_]: Applicative]: HelloWorld[F] = new HelloWorld[F]:
    def hello(n: HelloWorld.Name): F[HelloWorld.Greeting] =
      Greeting("Hello, " + n.name).pure[F]  // Create a Greeting by appending "Hello, " to the provided name
