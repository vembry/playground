package vembry.appscala

import cats.effect.{IO, IOApp}

object Main extends IOApp.Simple:
  // Entry point for the application, which runs the server
  val run = AppscalaServer.run[IO]
