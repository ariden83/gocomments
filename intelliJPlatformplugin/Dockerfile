# Étape 1 : Utiliser une image de base avec JDK 11
FROM openjdk:11-jdk

ENV GRADLE_VERSION=7.4.2

# Étape 4 : Installer Gradle
RUN wget https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip \
    && unzip gradle-${GRADLE_VERSION}-bin.zip -d /opt/ \
    && rm gradle-${GRADLE_VERSION}-bin.zip \
    && ln -s /opt/gradle-${GRADLE_VERSION} /opt/gradle

ENV PATH="/opt/gradle/bin:/opt/idea-IC-${IDEA_VERSION}/bin:${PATH}"

# Étape 5 : Configurer le répertoire de travail
WORKDIR /workspace

# Copier le fichier d'entrée (plugin.xml) dans le conteneur
COPY plugin.xml /workspace/plugin.xml

# Copiez le code source du plugin
COPY src /workspace/src
