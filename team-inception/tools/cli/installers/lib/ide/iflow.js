const path = require('node:path');
const { BaseIdeSetup } = require('./_base-ide');
const chalk = require('chalk');
const { AgentCommandGenerator } = require('./shared/agent-command-generator');

/**
 * iFlow CLI setup handler
 * Creates commands in .iflow/commands/ directory structure
 */
class IFlowSetup extends BaseIdeSetup {
  constructor() {
    super('iflow', 'iFlow CLI');
    this.configDir = '.iflow';
    this.commandsDir = 'commands';
  }

  /**
   * Setup iFlow CLI configuration
   * @param {string} projectDir - Project directory
   * @param {string} airDir - AIRCHETIPO installation directory
   * @param {Object} options - Setup options
   */
  async setup(projectDir, airDir, options = {}) {
    console.log(chalk.cyan(`Setting up ${this.name}...`));

    // Create .iflow/commands/air directory structure
    const iflowDir = path.join(projectDir, this.configDir);
    const commandsDir = path.join(iflowDir, this.commandsDir, 'air');
    const agentsDir = path.join(commandsDir, 'agents');
    const tasksDir = path.join(commandsDir, 'tasks');

    await this.ensureDir(agentsDir);
    await this.ensureDir(tasksDir);

    // Generate agent launchers
    const agentGen = new AgentCommandGenerator(this.airFolderName);
    const { artifacts: agentArtifacts } = await agentGen.collectAgentArtifacts(airDir, options.selectedModules || []);

    // Setup agents as commands
    let agentCount = 0;
    for (const artifact of agentArtifacts) {
      const commandContent = await this.createAgentCommand(artifact);

      const targetPath = path.join(agentsDir, `${artifact.module}-${artifact.name}.md`);
      await this.writeFile(targetPath, commandContent);
      agentCount++;
    }

    // Get tasks
    const tasks = await this.getTasks(airDir);

    // Setup tasks as commands
    let taskCount = 0;
    for (const task of tasks) {
      const content = await this.readFile(task.path);
      const commandContent = this.createTaskCommand(task, content);

      const targetPath = path.join(tasksDir, `${task.module}-${task.name}.md`);
      await this.writeFile(targetPath, commandContent);
      taskCount++;
    }

    console.log(chalk.green(`✓ ${this.name} configured:`));
    console.log(chalk.dim(`  - ${agentCount} agent commands created`));
    console.log(chalk.dim(`  - ${taskCount} task commands created`));
    console.log(chalk.dim(`  - Commands directory: ${path.relative(projectDir, commandsDir)}`));

    return {
      success: true,
      agents: agentCount,
      tasks: taskCount,
    };
  }

  /**
   * Create agent command content
   */
  async createAgentCommand(artifact) {
    // The launcher content is already complete - just return it as-is
    return artifact.content;
  }

  /**
   * Create task command content
   */
  createTaskCommand(task, content) {
    // Extract task name
    const nameMatch = content.match(/<name>([^<]+)<\/name>/);
    const taskName = nameMatch ? nameMatch[1] : this.formatTitle(task.name);

    let commandContent = `# /task-${task.name} Command

When this command is used, execute the following task:

## ${taskName} Task

${content}

## Usage

This command executes the ${taskName} task from the AIRCHETIPO ${task.module.toUpperCase()} module.

## Module

Part of the AIRCHETIPO ${task.module.toUpperCase()} module.
`;

    return commandContent;
  }

  /**
   * Cleanup iFlow configuration
   */
  async cleanup(projectDir) {
    const fs = require('fs-extra');
    const airCommandsDir = path.join(projectDir, this.configDir, this.commandsDir, 'air');

    if (await fs.pathExists(airCommandsDir)) {
      await fs.remove(airCommandsDir);
      console.log(chalk.dim(`Removed AIRCHETIPO commands from iFlow CLI`));
    }
  }
}

module.exports = { IFlowSetup };
