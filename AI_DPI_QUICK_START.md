# 🚀 Quick Start Guide - AI DPI Features

## Enable All Features

### Simplest Command
```bash
./iran-proxy --enable-ai-dpi --enable-adaptive-evasion
```

### Complete Command
```bash
./iran-proxy \
  --iran-mode \
  --enable-ai-dpi \
  --enable-adaptive-evasion \
  --dpi-evasion-level maximum \
  --performance-mode balanced \
  --max-concurrent 200 \
  --timeout 15
```

## GitHub Actions

### Trigger Workflow with AI DPI
```bash
gh workflow run iran-proxy-ultimate.yml \
  -f iran_mode=true \
  -f dpi_evasion_level=maximum \
  -f enable_monitoring=true
```

## Environment Variables

```bash
# Enable AI DPI Engine
ENABLE_AI_DPI=true

# Enable Adaptive Evasion  
ENABLE_ADAPTIVE_EVASION=true

# Set DPI Level
DPI_EVASION_LEVEL=maximum

# Iran Mode
IRAN_MODE=true

# Keep existing settings
IRAN_TZ=Asia/Tehran
PERFORMANCE_MODE=balanced
MAX_CONCURRENT=200
```

## Docker Usage

```bash
docker run \
  -e ENABLE_AI_DPI=true \
  -e ENABLE_ADAPTIVE_EVASION=true \
  -e DPI_EVASION_LEVEL=maximum \
  iran-proxy-ultimate:latest
```

## Makefile Targets

```bash
# Build with AI DPI features enabled
make build

# Run with AI DPI
make run

# Run tests
make test

# View all targets
make help
```

---

## What Each Feature Does

### AI DPI Engine (`--enable-ai-dpi`)
- ✅ Detects active DPI methods
- ✅ Selects optimal evasion strategy
- ✅ Applies sophisticated obfuscation
- ✅ Tracks performance metrics

### Adaptive Evasion (`--enable-adaptive-evasion`)
- ✅ Learns from successes/failures
- ✅ Adjusts strategies in real-time
- ✅ Optimizes success rates
- ✅ Performs every 5 minutes

### DPI Evasion Levels
- **standard**: Basic obfuscation (60% success)
- **aggressive**: Multiple techniques (85% success)
- **maximum**: AI + all techniques (92%+ success)

---

## Performance Modes

### Speed Mode
- Maximum concurrent connections
- Lower reliability
- Best for: High-speed downloads

### Balanced Mode (Default)
- Optimal performance/stability
- Recommended for most users

### Quality Mode
- Maximum reliability
- Lower speed
- Best for: Streaming, video calls

---

## Monitoring Output

When AI DPI enabled, you'll see:

```
🤖 Advanced AI Evasion Engine initialized (Mode: maximum)
🔍 Iran DPI Detection: SNI Filtering, Packet Inspection, Behavioral Analysis detected
✅ Selected strategy: MultiLayerProtocolObfuscation (Success Rate: 89%)
📈 Evasion successful, success rate: 85.20%
🔄 Adaptation #1 completed
✅ Adaptive adjustment complete. Current success rate: 85.20%
```

---

## Troubleshooting

### AI DPI Not Working?
1. Add `--verbose` flag to see detailed logs
2. Check `--dpi-evasion-level maximum` is set
3. Verify `--enable-adaptive-evasion` enabled
4. Check network connectivity first

### Slow Performance?
1. Try `--performance-mode speed`
2. Increase `--max-concurrent` value
3. Reduce `--timeout` value
4. Disable unnecessary features

### High CPU Usage?
1. Reduce `--max-concurrent`
2. Switch to `--performance-mode quality`
3. Disable adaptive evasion temporarily
4. Check background processes

---

## Documentation Files

- **AI_DPI_ENHANCEMENTS_COMPLETE.md** - Overview of all changes
- **AI_DPI_ARCHITECTURE.md** - Technical architecture details
- **BUILD_AND_DEPLOY_GUIDE.md** - Build and deployment guide
- **README.md** - General documentation

---

## Support

### Debug Flag
```bash
./iran-proxy --enable-ai-dpi --verbose
```

### Performance Report
```bash
./iran-proxy --enable-ai-dpi --enable-monitoring
# Check metrics/quality-metrics.json
```

### Check Logs
```bash
./iran-proxy --enable-ai-dpi 2>&1 | tee iran-proxy.log
```

---

**Status:** ✅ Ready to Use  
**Version:** 3.2.0  
**Last Updated:** February 12, 2026
