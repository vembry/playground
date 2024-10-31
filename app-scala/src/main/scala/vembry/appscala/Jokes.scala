package vembry.appscala

import cats.effect.Concurrent
import cats.syntax.all.*
import io.circe.{Encoder, Decoder}
import org.http4s.*
import org.http4s.implicits.*
import org.http4s.client.Client
import org.http4s.client.dsl.Http4sClientDsl
import org.http4s.circe.*
import org.http4s.Method.*

trait Jokes[F[_]]:
  // Defines a method to get a joke
  def get: F[Jokes.Joke]

object Jokes:
  // Apply method to summon an instance of Jokes
  def apply[F[_]](using ev: Jokes[F]): Jokes[F] = ev

  // Data structure for a joke response
  final case class Joke(joke: String)
  object Joke:
    // JSON decoder to convert JSON data into a Joke instance
    given Decoder[Joke] = Decoder.derived[Joke]

    // EntityDecoder for deserializing Joke from an HTTP response body
    given [F[_]: Concurrent]: EntityDecoder[F, Joke] = jsonOf

    // JSON encoder to convert Joke into JSON format
    given Encoder[Joke] = Encoder.AsObject.derived[Joke]

    // EntityEncoder for serializing Joke to an HTTP response body
    given [F[_]]: EntityEncoder[F, Joke] = jsonEncoderOf

  // Custom exception to handle joke-fetching errors
  final case class JokeError(e: Throwable) extends RuntimeException

  // Implementation of Jokes service to fetch a joke from an external API
  def impl[F[_]: Concurrent](C: Client[F]): Jokes[F] = new Jokes[F]:
    val dsl = new Http4sClientDsl[F]{}  // HTTP client DSL for request building
    import dsl.*
    def get: F[Jokes.Joke] = 
      C.expect[Joke](GET(uri"https://icanhazdadjoke.com/")) // Fetch joke as JSON
        .adaptError{ case t => JokeError(t) } // Catch any errors during the HTTP request and decoding
