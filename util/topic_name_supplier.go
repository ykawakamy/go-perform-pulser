package util

import (
	"regexp"
	"strconv"
)

// Topicの払い出しを行う機能を提供します。
// 原則1クライアントに対して、１つ生成する想定です。
type TopicNameSupplier struct {
	seq    uint
	topics []string
}

func CreateTopicNameSupplier() TopicNameSupplier {
	return TopicNameSupplier{
		seq:    0,
		topics: []string{"test", "test2"},
	}
}

// Topicの払い出しを行う。最後まで払いだしたら、初めから払いだす。
func (s *TopicNameSupplier) Get() string {
	ret := s.topics[int(s.seq)%len(s.topics)]
	s.seq += 1

	return ret
}

// 払いだすすべてのTopicを返す
func (s *TopicNameSupplier) GetAll() []string {
	return s.topics
}

// TopicNameSupplierを生成する機能を提供します。
type TopicNameSupplierFactory struct {
	topicPrefix     string
	topicSuffixExpr string
	topics          []string

	distoribution int
	current       int
}

func CreateFactory() *TopicNameSupplierFactory {
	return &TopicNameSupplierFactory{
		distoribution: 1,
	}
}

// TODO 説明を書く
func (s *TopicNameSupplierFactory) ParseTopicExpression(expr string) *TopicNameSupplierFactory {
	if len(s.topics) > 0 {
		panic("")
	}
	pattern := regexp.MustCompile(`([^:]*)(:(.+))?`)
	token := pattern.FindStringSubmatch(expr)
	if token == nil {
		return s
	}
	s.topicPrefix = token[1]
	s.topicSuffixExpr = token[3]

	if s.topicSuffixExpr == "" {
		s.topics = append(s.topics, s.topicPrefix)
		return s
	}

	rangePattern := regexp.MustCompile(`(\d+)-(\d+)`)
	rangeToken := rangePattern.FindStringSubmatch(s.topicSuffixExpr)
	if rangeToken != nil {
		start, _ := strconv.Atoi(rangeToken[1])
		end, _ := strconv.Atoi(rangeToken[2])
		for i := start; i <= end; i++ {
			tmp := s.topicPrefix + strconv.Itoa(i)
			s.topics = append(s.topics, tmp)
		}
		return s
	}

	tuplePattern := regexp.MustCompile(`,+`)
	tupleToken := tuplePattern.Split(s.topicSuffixExpr, -1)
	if tupleToken != nil {
		s.topics = append(s.topics, tupleToken...)
	}

	return s
}

func (s *TopicNameSupplierFactory) SetDistoribution(v int) *TopicNameSupplierFactory {
	if v <= 0 {
		panic("invalid distoribution value < 1")
	}
	s.distoribution = v

	return s
}

func (s *TopicNameSupplierFactory) Build() *TopicNameSupplier {
	length := len(s.topics)
	start := s.current * length / s.distoribution
	end := (s.current + 1) * length / s.distoribution
	s.current++
	if s.current == s.distoribution {
		end = length
		s.current = 0
	}
	part_topics := s.topics[start:end]

	return &TopicNameSupplier{
		seq:    0,
		topics: part_topics,
	}
}
