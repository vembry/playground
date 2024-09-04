# Use the official Gradle image to build the application
FROM eclipse-temurin:21-alpine AS build-stage

# Set the working directory
WORKDIR /home/app

# Copy Gradle files
COPY ./app-java/gradle ./gradle
COPY ./app-java/gradlew ./
COPY ./app-java/gradlew.bat ./
COPY ./app-java/build.gradle ./
COPY ./app-java/settings.gradle ./

# Copy source code
COPY ./app-java/src/main ./src/main

# Copy properties file
COPY ./app-java/src/main/resources/application.properties.docker ./src/main/resources/application.properties

# Run Gradle build
RUN ./gradlew build --no-daemon

# release stage
FROM gcr.io/distroless/java21-debian12 AS build-release-stage

# Set the working directory
WORKDIR /home/app

# Copy the built JAR file from the previous stage
COPY --from=build-stage /home/app/build/libs/*.jar /home/app/app.jar

# Expose the port the application runs on
EXPOSE 80

# Run the application
ENTRYPOINT ["java", "-jar", "/home/app/app.jar"]
