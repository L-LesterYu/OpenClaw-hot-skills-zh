package promptcatalog

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/geekjourneyx/md2wechat-skill/internal/assets"
	"gopkg.in/yaml.v3"
)

const promptsDirEnvVar = "MD2WECHAT_PROMPTS_DIR"

type PromptSpec struct {
	Name                    string            `yaml:"name" json:"name"`
	Kind                    string            `yaml:"kind" json:"kind"`
	Description             string            `yaml:"description" json:"description"`
	Version                 string            `yaml:"version" json:"version"`
	Archetype               string            `yaml:"archetype,omitempty" json:"archetype,omitempty"`
	PrimaryUseCase          string            `yaml:"primary_use_case,omitempty" json:"primary_use_case,omitempty"`
	CompatibleUseCases      []string          `yaml:"compatible_use_cases,omitempty" json:"compatible_use_cases,omitempty"`
	RecommendedAspectRatios []string          `yaml:"recommended_aspect_ratios,omitempty" json:"recommended_aspect_ratios,omitempty"`
	DefaultAspectRatio      string            `yaml:"default_aspect_ratio,omitempty" json:"default_aspect_ratio,omitempty"`
	Tags                    []string          `yaml:"tags,omitempty" json:"tags,omitempty"`
	Examples                []string          `yaml:"examples,omitempty" json:"examples,omitempty"`
	Variables               []string          `yaml:"variables,omitempty" json:"variables,omitempty"`
	Template                string            `yaml:"template" json:"template"`
	Metadata                map[string]string `yaml:"metadata,omitempty" json:"metadata,omitempty"`
	Source                  string            `yaml:"-" json:"source,omitempty"`
}

type ListFilter struct {
	Kind      string
	Archetype string
	Tag       string
}

type Catalog struct {
	prompts map[string]*PromptSpec
}

var (
	defaultCatalogOnce sync.Once
	defaultCatalog     *Catalog
	defaultCatalogErr  error
)

func NewCatalog() *Catalog {
	return &Catalog{prompts: make(map[string]*PromptSpec)}
}

func DefaultCatalog() (*Catalog, error) {
	defaultCatalogOnce.Do(func() {
		cat := NewCatalog()
		defaultCatalogErr = cat.Load()
		if defaultCatalogErr == nil {
			defaultCatalog = cat
		}
	})
	return defaultCatalog, defaultCatalogErr
}

func ResetDefaultCatalogForTests() {
	defaultCatalogOnce = sync.Once{}
	defaultCatalog = nil
	defaultCatalogErr = nil
}

func key(kind, name string) string {
	return strings.ToLower(strings.TrimSpace(kind)) + "/" + strings.ToLower(strings.TrimSpace(name))
}

func (c *Catalog) Load() error {
	if err := c.loadBuiltin(); err != nil {
		return err
	}
	for _, dir := range getPromptDirs() {
		if err := c.loadFromDir(dir); err != nil {
			return err
		}
	}
	return nil
}

func getPromptDirs() []string {
	homeDir, _ := os.UserHomeDir()
	dirs := []string{
		filepath.Join(homeDir, ".config", "md2wechat", "prompts"),
		"prompts",
	}
	if explicitDir := strings.TrimSpace(os.Getenv(promptsDirEnvVar)); explicitDir != "" {
		dirs = append(dirs, explicitDir)
	}
	return dirs
}

func (c *Catalog) loadBuiltin() error {
	kinds := []string{"humanizer", "refine", "image"}
	for _, kind := range kinds {
		names, err := assets.ListBuiltinPrompts(kind)
		if err != nil {
			return fmt.Errorf("list builtin prompts for %s: %w", kind, err)
		}
		for _, name := range names {
			data, err := assets.ReadBuiltinPrompt(kind, name)
			if err != nil {
				return fmt.Errorf("read builtin prompt %s/%s: %w", kind, name, err)
			}
			spec, err := parsePromptSpec(data)
			if err != nil {
				return fmt.Errorf("parse builtin prompt %s/%s: %w", kind, name, err)
			}
			spec.Source = "builtin"
			c.prompts[key(spec.Kind, spec.Name)] = spec
		}
	}
	return nil
}

func (c *Catalog) loadFromDir(root string) error {
	if root == "" {
		return nil
	}
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read prompt directory %s: %w", root, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		kindDir := filepath.Join(root, entry.Name())
		files, err := os.ReadDir(kindDir)
		if err != nil {
			return fmt.Errorf("read prompt kind directory %s: %w", kindDir, err)
		}
		for _, file := range files {
			if file.IsDir() {
				continue
			}
			if !strings.HasSuffix(file.Name(), ".yaml") && !strings.HasSuffix(file.Name(), ".yml") {
				continue
			}
			data, err := os.ReadFile(filepath.Join(kindDir, file.Name()))
			if err != nil {
				return err
			}
			spec, err := parsePromptSpec(data)
			if err != nil {
				return err
			}
			spec.Source = root
			c.prompts[key(spec.Kind, spec.Name)] = spec
		}
	}
	return nil
}

func parsePromptSpec(data []byte) (*PromptSpec, error) {
	var spec PromptSpec
	if err := yaml.Unmarshal(data, &spec); err != nil {
		return nil, err
	}
	if spec.Kind == "" || spec.Name == "" {
		return nil, fmt.Errorf("prompt spec requires kind and name")
	}
	if spec.Template == "" {
		return nil, fmt.Errorf("prompt spec requires template")
	}
	if spec.PrimaryUseCase == "" && spec.Archetype != "" {
		spec.PrimaryUseCase = spec.Archetype
	}
	return &spec, nil
}

func (c *Catalog) List(kind string) []PromptSpec {
	return c.ListFiltered(ListFilter{Kind: kind})
}

func (c *Catalog) ListFiltered(filter ListFilter) []PromptSpec {
	result := make([]PromptSpec, 0, len(c.prompts))
	for _, spec := range c.prompts {
		if filter.Kind != "" && !strings.EqualFold(spec.Kind, filter.Kind) {
			continue
		}
		if filter.Archetype != "" && !strings.EqualFold(spec.Archetype, filter.Archetype) {
			continue
		}
		if filter.Tag != "" && !containsTag(spec.Tags, filter.Tag) {
			continue
		}
		result = append(result, *spec)
	}
	sort.Slice(result, func(i, j int) bool {
		if result[i].Kind == result[j].Kind {
			return result[i].Name < result[j].Name
		}
		return result[i].Kind < result[j].Kind
	})
	return result
}

func containsTag(tags []string, target string) bool {
	target = strings.TrimSpace(target)
	if target == "" {
		return true
	}
	for _, tag := range tags {
		if strings.EqualFold(tag, target) {
			return true
		}
	}
	return false
}

func SupportsUseCase(spec *PromptSpec, useCase string) bool {
	if spec == nil {
		return false
	}
	useCase = strings.TrimSpace(useCase)
	if useCase == "" {
		return true
	}
	if strings.EqualFold(spec.PrimaryUseCase, useCase) {
		return true
	}
	if strings.EqualFold(spec.Archetype, useCase) {
		return true
	}
	for _, compatible := range spec.CompatibleUseCases {
		if strings.EqualFold(compatible, useCase) {
			return true
		}
	}
	return false
}

func (c *Catalog) Get(kind, name string) (*PromptSpec, error) {
	if kind != "" {
		spec, ok := c.prompts[key(kind, name)]
		if !ok {
			return nil, fmt.Errorf("prompt not found: %s/%s", kind, name)
		}
		copy := *spec
		return &copy, nil
	}

	var match *PromptSpec
	for _, spec := range c.prompts {
		if strings.EqualFold(spec.Name, name) {
			if match != nil {
				return nil, fmt.Errorf("prompt name %s is ambiguous; specify kind", name)
			}
			copy := *spec
			match = &copy
		}
	}
	if match == nil {
		return nil, fmt.Errorf("prompt not found: %s", name)
	}
	return match, nil
}

func (c *Catalog) Render(kind, name string, vars map[string]string) (string, *PromptSpec, error) {
	spec, err := c.Get(kind, name)
	if err != nil {
		return "", nil, err
	}
	rendered := spec.Template
	for key, value := range vars {
		rendered = strings.ReplaceAll(rendered, "{{"+key+"}}", value)
	}
	return rendered, spec, nil
}
