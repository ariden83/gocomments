package com.ariden.gocomments;

import com.intellij.AppTopics;
import com.intellij.openapi.application.ApplicationManager;
import com.intellij.openapi.command.WriteCommandAction;
import com.intellij.openapi.components.ApplicationComponent;
import com.intellij.openapi.editor.Document;
import com.intellij.openapi.fileEditor.FileDocumentManager;
import com.intellij.openapi.fileEditor.FileDocumentManagerListener;
import com.intellij.openapi.project.Project;
import com.intellij.openapi.ui.Messages;
import com.intellij.openapi.vfs.VirtualFile;
import com.intellij.util.messages.MessageBusConnection;
import org.jetbrains.annotations.NotNull;
import com.intellij.openapi.diagnostic.Logger;

import java.io.BufferedReader;
import java.io.InputStreamReader;

public class SaveActionListener implements ApplicationComponent {

    private static final Logger LOG = Logger.getInstance(SaveActionListener.class);

    @Override
    public void initComponent() {
        LOG.info("Initializing SaveActionListener component");
        MessageBusConnection connection = ApplicationManager.getApplication().getMessageBus().connect();
        connection.subscribe(AppTopics.FILE_DOCUMENT_SYNC, new FileDocumentManagerListener() {
            @Override
            public void beforeDocumentSaving(@NotNull Document document) {
                VirtualFile file = FileDocumentManager.getInstance().getFile(document);
                if (file != null && file.getFileType().getName().equals("Go")) {
                    Project project = ProjectLocator.getInstance().guessProjectForFile(file);
                    if (project != null) {
                        LOG.info("Before saving document: " + file.getPath());
                        runGoComments(project, file);
                    }
                }
            }
        });
    }

    private void runGoComments(Project project, VirtualFile file) {
        LOG.info("Running GoComments on file: " + file.getPath());
        ApplicationManager.getApplication().invokeLater(() -> WriteCommandAction.runWriteCommandAction(project, () -> {
            try {
                // Chemin vers le binaire compilé
                String binaryPath = "./bin/gocomments"; // ajustez le chemin selon votre configuration

                // Récupérer le nom du fichier et son chemin
                String filePath = file.getPath();
                String fileName = file.getName();

                // Créer la commande à exécuter
                ProcessBuilder processBuilder = new ProcessBuilder(binaryPath, "-l", "-w", filePath);
                processBuilder.redirectErrorStream(true);

                // Exécuter la commande
                Process process = processBuilder.start();

                // Lire la sortie de la commande
                BufferedReader reader = new BufferedReader(new InputStreamReader(process.getInputStream()));
                StringBuilder output = new StringBuilder();
                String line;
                while ((line = reader.readLine()) != null) {
                    output.append(line).append("\n");
                }

                // Attendre la fin de l'exécution de la commande
                int exitCode = process.waitFor();

                // Afficher le résultat si nécessaire
                if (exitCode != 0) {
                    LOG.error("Error executing gocomments on " + fileName + ": Exit code " + exitCode);
                    Messages.showMessageDialog(project, "Error executing gocomments on " + fileName, "Error", Messages.getErrorIcon());
                } else {
                    LOG.info("gocomments executed successfully on " + fileName);
                }
            } catch (Exception ex) {
                LOG.error("Exception while executing gocomments on " + file.getPath(), ex);
                ex.printStackTrace();
            }
        }));
    }

    @Override
    public void disposeComponent() {
        LOG.info("Disposing SaveActionListener component");
    }

    @NotNull
    @Override
    public String getComponentName() {
        return "SaveActionListener";
    }
}
