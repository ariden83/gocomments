# gocomments

![languague: Go](https://img.shields.io/badge/language-go-007d9c)

`gocomments` is a revolutionary Go documentation tool that creates and trains **custom AI models** to automatically generate intelligent comments for Go functions, variables, and types.

## Custom AI Model Training System

This project implements a complete machine learning pipeline that trains **your own AI models** specifically for Go code documentation:

### What Makes This Unique

- **Custom T5-Based Model**: Built on Salesforce's CodeT5 foundation model, fine-tuned for Go documentation
- **Self-Supervised Learning**: Trains from your existing Go codebase comments to learn your coding patterns
- **TensorFlow Implementation**: Complete neural network training pipeline with GPU support
- **Incremental Learning**: Models improve through iterative training with quality scoring
- **REST API Deployment**: Trained models are deployed as microservices for real-time inference

### Training Process

1. **Data Extraction**: Parses Go source files to extract function signatures and existing comments
2. **Quality Analysis**: Uses spaCy NLP to analyze comment quality and filter training data
3. **Model Training**: Fine-tunes CodeT5 transformer model using TensorFlow
4. **Evaluation**: Tests multiple model versions to find the best performer
5. **Deployment**: Serves trained models via Flask API for real-time comment generation


## Quick Start

### 1. Install the Tool

```bash
make install
```

### 2. Train Your Custom AI Model

#### Generate Training Dataset from Your Codebase

Create a `.env` file in the `dataset/` folder:

```bash
LOCAL_REPO_PATH=/path/to/your/go/repository
```

Extract function-comment pairs from your codebase:

```bash
make generate-dataset
```

#### Train the Neural Network

Launch the TensorFlow training pipeline:

```bash
make generate-model
```

This will:
- Fine-tune a CodeT5 transformer model on your data
- Train for 60 epochs with validation
- Save multiple model checkpoints
- Generate training/validation loss curves

#### Deploy the Trained Model

Start the inference API server:

```bash
make generate-api
```

Your trained model is now available at `http://localhost:5000`

### 3. Generate Comments with Your AI

Configure your project with a `.gocomments` file:

```yaml
localai:
  active: true
  url: "http://localhost:5000"
  api_model_version: 10  # Use the best performing model version
```

Generate comments using your trained AI:

```bash
gocomments -l -w ./src/.
```

## Advanced Usage

### Command Line Options

```text
Usage: gocomments [flags] [path ...]
  -d	display diffs instead of rewriting files
  -l	list files whose formatting differs from goimport's
  -local string
    	put imports beginning with this string after 3rd-party package
  -prefix value
    	relative local prefix to from a new import group (can be given several times)
  -w	write result to (source) file instead of stdout
```

### Testing Model Performance

Evaluate different model versions:

```bash
make generate-test
```

This compares all trained model versions against test data to find the best performer.

### Configuration File Format

Create a `.gocomments` file in your project directory:

```yaml
---
# Basic configuration
local: ""  # Root Go module name (auto-detected from go.mod)
prefixes:
  - "common"  # Import grouping prefixes
signature: "AutoComBOT"  # Comment signature for tracking
update-comments: false  # Update existing AI-generated comments
active-examples: true   # Generate usage examples in comments

# Your Custom AI Model Configuration
localai:
  active: true
  url: "http://localhost:5000"
  api_model_version: 10  # Specify which trained model version to use
```

## Deep Dive: AI Model Architecture

### Neural Network Details

- **Base Model**: Salesforce CodeT5 (Text-to-Text Transfer Transformer)
- **Architecture**: Encoder-Decoder transformer with 220M parameters
- **Training Framework**: TensorFlow with mixed precision support
- **Tokenizer**: RoBERTa-based tokenizer optimized for code
- **Training Strategy**: Fine-tuning with distributed training support

### Data Processing Pipeline

#### 1. Code Analysis (`parse_functions.go`)
```go
type FunctionInfo struct {
    Name    string `json:"name"`
    Comment string `json:"comment"`
}
```

- Parses Go AST to extract function signatures
- Filters out test functions, `main()`, and `init()`
- Excludes comments with keywords like "TODO", "FIXME", "deprecated"

#### 2. Quality Scoring (`generate_func_comments_from_local_repo.py`)
```python
def evaluate_quality(doc):
    # Analyzes comment quality using spaCy NLP:
    # - Minimum length and verb requirements
    # - Technical term density
    # - Stop word ratio analysis
    return "good" | "average" | "poor"
```

#### 3. Model Training (`train.py`)
- **Input**: `{file_name} {function_signature}`
- **Output**: `{comment} POS: {pos_tags}` (includes grammatical analysis)
- **Training**: 60 epochs with validation split
- **Optimization**: AdamW optimizer with learning rate scheduling

### Model Performance Analysis

The system tracks multiple model versions and compares their performance:

```
Model Version 1:  Basic function understanding
Model Version 5:  Improved context awareness  
Model Version 10: Domain-specific optimization
```

Each version is evaluated on:
- **Semantic Accuracy**: How well comments describe function purpose
- **Syntactic Quality**: Grammar and technical terminology usage
- **Consistency**: Alignment with existing codebase patterns

## AI-Generated Documentation Examples

## Technical Architecture

### System Components

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Data Layer    │    │  Training Layer │    │ Inference Layer │
│                 │    │                 │    │                 │
│ • Go AST Parser │───▶│ • TensorFlow    │───▶│ • Flask API     │
│ • spaCy NLP     │    │ • CodeT5 Model  │    │ • Model Serving │
│ • Quality Filter│    │ • GPU Training  │    │ • Go Integration│
└─────────────────┘    └─────────────────┘    └─────────────────┘
```

### Directory Structure

```
gocomments/
├── dataset/                    # Data processing & extraction
│   ├── generate_func_comments_from_local_repo.py
│   ├── parse_functions.go
│   └── file/                   # Generated training datasets
├── model/                      # Neural network training
│   ├── train.py               # TensorFlow training pipeline
│   ├── Dockerfile             # Containerized training
│   └── runs/                  # Model checkpoints & logs
├── api/                       # Model inference server
│   ├── tokenizer_api.py       # Flask API for trained models
│   └── docker-compose.yml     # API deployment
├── test-models/               # Model evaluation
│   └── main.go               # Performance testing
└── internal/comments/         # Go tool integration
    ├── comments_localai.go    # Custom AI integration
    └── comments_interface.go  # Provider interface
```

## Future Enhancements

### Planned Features

- **Multi-Language Support**: Extend to Python, JavaScript, Rust
- **Context-Aware Comments**: Analyze broader code context for better suggestions
- **IDE Integrations**: VS Code, GoLand, and Vim plugins
- **Continuous Learning**: Model updates from user feedback
- **Enterprise Features**: Team model sharing and version control

### Research Areas

- **Code Understanding**: Improve semantic analysis of complex functions
- **Comment Templates**: Generate structured documentation formats
- **Performance Optimization**: Faster inference and smaller model sizes

## Research & References

### Foundation Models
- [CodeT5: Identifier-aware Unified Pre-trained Encoder-Decoder Models](https://arxiv.org/abs/2109.00859)
- [RoBERTa: A Robustly Optimized BERT Pretraining Approach](https://arxiv.org/abs/1907.11692)

### Code Generation Research
- [CodeSearchNet Challenge](https://arxiv.org/abs/1909.09436) - Code search and documentation
- [MBPP: Mostly Basic Python Problems](https://arxiv.org/abs/2108.07732) - Code generation benchmark

### NLP for Code
- [spaCy Industrial-Strength Natural Language Processing](https://spacy.io/)
- [Transformers: State-of-the-Art Natural Language Processing](https://arxiv.org/abs/1910.03771)

## Disclaimer

This tool generates AI-powered documentation. Always review generated comments for accuracy and appropriateness before committing to production code.
