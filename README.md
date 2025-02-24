# Pocket Doc - AI-Powered Medical Assistant
<img src="https://github.com/Boot41/Itish-fse/blob/main/client/doctor-ai/public/Screenshot%20from%202025-02-24%2014-24-58.png"/>

## Overview
Pocket Doc is an innovative platform that leverages artificial intelligence to assist medical professionals in their daily tasks. The platform provides features like AI-powered voice transcription, automated report generation, and advanced analytics.

## Features
- **AI-Powered Transcription**: Real-time voice-to-text conversion with medical terminology recognition
- **Automated Reporting**: Generate comprehensive medical reports from transcriptions
- **Patient Management**: Track and manage patient information and history
- **Analytics Dashboard**: Gain insights into consultation trends and practice performance
- **Secure Authentication**: HIPAA-compliant security measures
- **Documentation & Support**: Comprehensive guides and support resources

## Tech Stack
### Frontend
- React with TypeScript
- Tailwind CSS for styling
- Vite for build tooling
- React Router for navigation
- Lucide React for icons

### Backend
- Go (Golang)
- Gin framework for routing
- JWT for authentication
- PostgreSQL for database
- AssemblyAI for speech-to-text
- Groq for text enchancement

## Installation

### Prerequisites
- Go 1.20+
- Node.js 18+
- PostgreSQL 14+
- Yarn or npm

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/doctor-ai.git
   cd doctor-ai
   ```

2. Install backend dependencies:
   ```bash
   cd doctor_ai
   go mod download
   ```

3. Install frontend dependencies:
   ```bash
   cd ../client/doctor-ai
   yarn install
   ```

4. Set up environment variables:
   Create `.env` files in both `doctor_ai` and `client/doctor-ai` directories with required configurations. Refer to `env_template` files in both `client/doctor-ai` and `doctor_ai`.

## Configuration
### Backend (Refer to `env_template` in `client/doctor-ai`)
Create a `.env` file in the `doctor_ai` directory:
```env
PORT=8080
JWT_SECRET=your-secret-key
DATABASE_URL=postgres://user:password@localhost:5432/doctor_ai (or connection string)
ASSEMBLYAI_KEY=your-assemblyai-key
```

### Frontend (Refer to `env_template` in `doctor_ai`)
Create a `.env` file in the `client/doctor-ai` directory:
```env
VITE_API_BASE_URL=http://localhost:8080
```

## Running the Project
1. Start the backend:
   ```bash
   cd doctor_ai
   go run main.go
   ```

2. Start the frontend:
   ```bash
   cd ../client/doctor-ai
   yarn dev
   ```

3. Access the application at `http://localhost:5173`

## Testing
First, install the required testing tools:
```bash
go install github.com/agiledragon/gomonkey/v2@latest
go install github.com/kyoh86/richgo@latest
export PATH=$PATH:$(go env GOPATH)/bin
source ~/.bashrc
```

### Normal Tests
```bash
# Verbose testing
richgo test -v ./...

# Count of tests
richgo test -v ./... 2>&1 | grep '^=== RUN' | wc -l
```

### Coverage Report
```bash
richgo test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## Deployment
The project can be deployed using Docker containers. Example `docker-compose.yml`:
```yaml
version: '3'
services:
  web:
    build: ./client/doctor-ai
    ports:
      - "5173:5173"
    depends_on:
      - api

  api:
    build: ./doctor_ai
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://user:password@db:5432/doctor_ai
    depends_on:
      - db

  db:
    image: postgres:14
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: doctor_ai
    ports:
      - "5432:5432"
```

## Upcoming Features (Pocket Doc V2)

### 1. Improved Patient Form
- Add an option to choose between uploading an audio file or providing a URL link.
- Enable image uploads for patient profile photos.
- Introduce a dropdown for patient type selection:
  - New Patient
  - Regular Patient
  - Returning Patient
  - Emergency Cases

### 2. Improved Dynamic Dashboard
- Modify the doctor model to include a field for storing data required for generating dynamic data.
- The backend logic for processing dynamic data is already implemented.

### 3. Improved Suggestions Chat Window
- Enhance the Suggestions Window to access and read data from transcripts.

## Contributing
Contributions are welcome! Please follow these steps:
1. Fork the repository
2. Create a new branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add some feature'`)
4. Push to the branch (`git push origin feature/YourFeature`)
5. Open a pull request

## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
