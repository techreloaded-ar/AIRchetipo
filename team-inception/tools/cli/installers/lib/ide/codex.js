const path = require('node:path');
const fs = require('fs-extra');
const os = require('node:os');
const chalk = require('chalk');
const { BaseIdeSetup } = require('./_base-ide');
const { WorkflowCommandGenerator } = require('./shared/workflow-command-generator');
const { AgentCommandGenerator } = require('./shared/agent-command-generator');
const { getTasksFromAir } = require('./shared/air-artifacts');

/**
 * Codex setup handler (CLI mode)
 */
class CodexSetup extends BaseIdeSetup {
  constructor() {
    super('codex', 'Codex', true); // preferred IDE
  }

  /**
   * Setup Codex configuration
   * @param {string} projectDir - Project directory
   * @param {string} airDir - AIRCHETIPO installation directory
   * @param {Object} options - Setup options
   */
  async setup(projectDir, airDir, options = {}) {
    console.log(chalk.cyan(`Setting up ${this.name}...`));

    // Always use CLI mode
    const mode = 'cli';

    const { artifacts, counts } = await this.collectClaudeArtifacts(projectDir, airDir, options);

    const destDir = this.getCodexPromptDir();
    await fs.ensureDir(destDir);
    await this.clearOldAIRchetipoFiles(destDir);
    const written = await this.flattenAndWriteArtifacts(artifacts, destDir);

    console.log(chalk.green(`✓ ${this.name} configured:`));
    console.log(chalk.dim(`  - Mode: CLI`));
    console.log(chalk.dim(`  - ${counts.agents} agents exported`));
    console.log(chalk.dim(`  - ${counts.tasks} tasks exported`));
    console.log(chalk.dim(`  - ${counts.workflows} workflow commands exported`));
    if (counts.workflowLaunchers > 0) {
      console.log(chalk.dim(`  - ${counts.workflowLaunchers} workflow launchers exported`));
    }
    console.log(chalk.dim(`  - ${written} Codex prompt files written`));
    console.log(chalk.dim(`  - Destination: ${destDir}`));

    // Prominent notice about home directory installation
    console.log('');
    console.log(chalk.bold.cyan('═'.repeat(70)));
    console.log(chalk.bold.yellow('  IMPORTANT: Codex Configuration'));
    console.log(chalk.bold.cyan('═'.repeat(70)));
    console.log('');
    console.log(chalk.white('  Prompts have been installed to your HOME DIRECTORY, not this project.'));
    console.log(chalk.white('  No .codex file was created in the project root.'));
    console.log('');
    console.log(chalk.green('  ✓ You can now use slash commands (/) in Codex CLI'));
    console.log(chalk.dim('    Example: /air-aim-agents-pm'));
    console.log(chalk.dim('    Type / to see all available commands'));
    console.log('');
    console.log(chalk.bold.cyan('═'.repeat(70)));
    console.log('');

    return {
      success: true,
      mode,
      artifacts,
      counts,
      destination: destDir,
      written,
    };
  }

  /**
   * Detect Codex installation by checking for AIRCHETIPO prompt exports
   */
  async detect(_projectDir) {
    const destDir = this.getCodexPromptDir();

    if (!(await fs.pathExists(destDir))) {
      return false;
    }

    const entries = await fs.readdir(destDir);
    return entries.some((entry) => entry.startsWith('air-'));
  }

  /**
   * Collect Claude-style artifacts for Codex export.
   * Returns the normalized artifact list for further processing.
   */
  async collectClaudeArtifacts(projectDir, airDir, options = {}) {
    const selectedModules = options.selectedModules || [];
    const artifacts = [];

    // Generate agent launchers
    const agentGen = new AgentCommandGenerator(this.airFolderName);
    const { artifacts: agentArtifacts } = await agentGen.collectAgentArtifacts(airDir, selectedModules);

    for (const artifact of agentArtifacts) {
      artifacts.push({
        type: 'agent',
        module: artifact.module,
        sourcePath: artifact.sourcePath,
        relativePath: artifact.relativePath,
        content: artifact.content,
      });
    }

    const tasks = await getTasksFromAir(airDir, selectedModules);
    for (const task of tasks) {
      const content = await this.readAndProcessWithProject(
        task.path,
        {
          module: task.module,
          name: task.name,
        },
        projectDir,
      );

      artifacts.push({
        type: 'task',
        module: task.module,
        sourcePath: task.path,
        relativePath: path.join(task.module, 'tasks', `${task.name}.md`),
        content,
      });
    }

    const workflowGenerator = new WorkflowCommandGenerator(this.airFolderName);
    const { artifacts: workflowArtifacts, counts: workflowCounts } = await workflowGenerator.collectWorkflowArtifacts(airDir);
    artifacts.push(...workflowArtifacts);

    return {
      artifacts,
      counts: {
        agents: agentArtifacts.length,
        tasks: tasks.length,
        workflows: workflowCounts.commands,
        workflowLaunchers: workflowCounts.launchers,
      },
    };
  }

  getCodexPromptDir() {
    return path.join(os.homedir(), '.codex', 'prompts');
  }

  flattenFilename(relativePath) {
    const sanitized = relativePath.replaceAll(/[\\/]/g, '-');
    return `air-${sanitized}`;
  }

  async flattenAndWriteArtifacts(artifacts, destDir) {
    let written = 0;

    for (const artifact of artifacts) {
      const flattenedName = this.flattenFilename(artifact.relativePath);
      const targetPath = path.join(destDir, flattenedName);
      await fs.writeFile(targetPath, artifact.content);
      written++;
    }

    return written;
  }

  async clearOldAIRchetipoFiles(destDir) {
    if (!(await fs.pathExists(destDir))) {
      return;
    }

    const entries = await fs.readdir(destDir);

    for (const entry of entries) {
      if (!entry.startsWith('air-')) {
        continue;
      }

      const entryPath = path.join(destDir, entry);
      const stat = await fs.stat(entryPath);
      if (stat.isFile()) {
        await fs.remove(entryPath);
      } else if (stat.isDirectory()) {
        await fs.remove(entryPath);
      }
    }
  }

  async readAndProcessWithProject(filePath, metadata, projectDir) {
    const content = await fs.readFile(filePath, 'utf8');
    return super.processContent(content, metadata, projectDir);
  }

  /**
   * Cleanup Codex configuration (no-op until export destination is finalized)
   */
  async cleanup() {
    const destDir = this.getCodexPromptDir();
    await this.clearOldAIRchetipoFiles(destDir);
  }
}

module.exports = { CodexSetup };
