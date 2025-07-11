# 👤 go-idvalidator

**go-idvalidator** es una librería escrita 100% en Go que permite analizar imágenes de documentos de identidad (DNI) y videos de rostros para validar la identidad de una persona. Diseñada como base para un sistema más robusto, esta librería busca ofrecer un primer nivel de coincidencia entre el documento físico y el rostro mediante procesamiento de imágenes.

---

## 🚀 Objetivo

Crear un wrapper/librería en Go que reciba como entrada:

- 2 imágenes del DNI (`frontal` y `posterior`)
- 1 video corto (o secuencia de imágenes) del rostro actual

Y entregue como salida:

```json
{
  "match_percentage": 92.5,
  "valid": true,
  "reason": "Coincidencia facial aceptable"
}
```

# 🧱 Tecnologías utilizadas
Go 1.21+

<!-- * Tesseract OCR para extraer texto del DNI -->
* go-face para análisis facial (bindings de Dlib)
* ffmpeg para extracción de frames desde el video
* Go image libraries para procesamiento

```bash
identity-validator/
├── go.mod
├── go.sum
├── README.md
├── cmd/
│   ├── demo/
│   │   └── main.go          # Ejecutable para realizar pruebas
│   └── main.go              # Ejecutable para el CLI
├── pkg/
│   ├── validator/
│   │   ├── analize.go       # Punto de entrada del análisis
│   │   └── types.go         # Definición de tipos Input/Result
├── internal/
│   ├── validation/
│   │   └── validation.go    # funcioenes base para el analisis
│   └── face/
│   │   ├── recognizer.go    # Comparacion de imagenes
│   │   └── models/          # Modelos descargados pre entrenados para el analisis
│   └── utils/
│   │   └── convertor.go     # Funciones para transformar png to jpg/jpeg
│   └── encryption/
│   │   └── encrypt.go       # Funciones para encriptar / desencriptar
│   └── config/
│       └── config.go        # Funciones para cargar secrets
└── testdata/
    ├── sample_front_dni.jpg
    ├── sample_back_dni.jpg
    └── sample_video_face.mp4
```

# 📥 Instalación
Requisitos previos
<!-- * Instalar Tesseract: -->

<!-- ```bash
brew install tesseract
```
* Instalar opencv:

```bash
brew upgrade opencv
brew install opencv
``` -->
* Instalar ffmpeg:
```bash
brew install ffmpeg jpeg libpng libtiff
```

* Dlib instalado (v19.10 o superior)
```bash
brew install dlib boost

# Modelos preentrenados de Dlib
wget http://dlib.net/files/shape_predictor_5_face_landmarks.dat.bz2
wget http://dlib.net/files/dlib_face_recognition_resnet_model_v1.dat.bz2
bunzip2 *.bz2
```


* Instalar librerías necesarias para go-face:

```bash
# Requiere CMake, Dlib y pkg-config
brew install cmake pkg-config
```

# ⚙️ Uso básico
Desde código Go

```go
input := validator.Input{
  FrontDNIPath: "testdata/dni_front.jpg",
  BackDNIPath:  "testdata/dni_back.jpg",
  VideoPath:    "testdata/face_video.mp4",
}

encrypted, err := validator.Analyze(input)
if err != nil {
  log.Fatal(err)
}

result, err := validator.DecryptAnalysis(encrypted)
if err != nil {
  log.Fatalf("Error in decrypting: %v", err)
}

fmt.Printf("Resultado: %+v\n", result)

```

Desde CLI
```bash
go run cmd/main.go \
  --front testdata/dni_front_1.png \
  --back testdata/dni_back.jpg \
  --video testdata/face_video.mp4 \
  --selfie testdata/selfie_0.png

```

# 📤 Respuesta esperada
```json
{
  "match_percentage": 94.3,
  "confidence_level": "high",
  "valid": true,
  "dni_text": {
    "name": "JUAN PEREZ",
    "id_number": "12345678-9"
  },
  "reason": "Coincidencia facial aceptable"
}
```

# 🧪 Roadmap inicial
- [x] OCR de cara frontal del DNI
- [x] Matching facial con go-face
- [x] Persistencia encriptada de resultados
- [ ] Extracción de frame facial del video
- [ ] Validación de estructura de DNI

# 🔐 Seguridad
Los datos procesados se consideran sensibles. En futuras versiones se incluirá:
* Encriptación AES de los resultados
* Firma digital de los resultados
* TTL de almacenamiento temporal

# 📄 Licencia
go-identity-validator is licensed under [CC0](LICENSE).
